package walker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	kafka "gopkg.in/Shopify/sarama.v1"

	"qiniu.com/app/common"
	log "qiniupkg.com/x/log.v7"

	"gopkg.in/mgo.v2"

	rpc "github.com/qiniu/rpc.v3"
	atypo "qiniu.com/app/avio-apiserver/typo"
	"qiniu.com/app/avio-walker/database"
	"qiniu.com/app/avio-walker/typo"
	"qiniu.com/app/avio-walker/utils"
)

type Walker interface {
	SetServerRouter(router *gin.Engine)
	Run() error
	AppendJob(job *atypo.JobInfo) error
	GracefullyShutdown() error

	// life cycle functions
	walkerWillStart() error
	walkerStarted() error
	walkerWillStop() error
	walkerStopped() error
}

type jobMessage struct {
	name   string
	status typo.WalkerStatus
}

const WALKER_CONCURRENT_LIMIT = 50

type walker struct {
	conf        *utils.Config
	name        string
	stopCh      chan bool
	mgoSession  *mgo.Session
	server      *http.Server
	jobChan     chan jobMessage
	wsUpgrader  *websocket.Upgrader
	walkCount   uint32
	wcWg        *sync.RWMutex
	rpcClient   *rpc.Client
	primaryHost string
	producer    kafka.SyncProducer
}

// TODO handle http server when life cycle function failed
func NewWalker(conf *utils.Config) (Walker, error) {
	w := &walker{}
	w.conf = conf
	session, e := database.NewMongoSession(conf)
	if e != nil {
		return nil, e
	}

	if e := database.Init(); e != nil {
		return nil, e
	}

	n, e := utils.GetWalkerName()
	if e != nil {
		return nil, e
	}
	w.name = n
	w.mgoSession = session
	w.stopCh = make(chan bool)
	w.jobChan = make(chan jobMessage)
	w.walkCount = 0
	w.wcWg = &sync.RWMutex{}
	w.rpcClient = &rpc.Client{
		Client: &http.Client{
			Timeout: time.Second * 300,
		},
	}

	w.primaryHost = "192.168.1.9:19999"

	// TODO handle web socket messages
	w.wsUpgrader = &websocket.Upgrader{}
	w.server = &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if e := w.initProducer(); e != nil {
		return nil, e
	}

	return w, nil
}

func (w *walker) createTLSConfiguration() (t *tls.Config) {
	if w.conf.KAFKA.CertFile != "" && w.conf.KAFKA.KeyFile != "" && w.conf.KAFKA.CaFile != "" {
		cert, err := tls.LoadX509KeyPair(w.conf.KAFKA.CertFile, w.conf.KAFKA.KeyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := ioutil.ReadFile(w.conf.KAFKA.CaFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: w.conf.KAFKA.VerifySsl,
		}
	}

	return t
}

func (w *walker) initProducer() error {
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := kafka.NewConfig()
	config.Producer.RequiredAcks = kafka.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                  // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	tlsConfig := w.createTLSConfiguration()
	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}

	producer, err := kafka.NewSyncProducer(w.conf.KAFKA.Brokers, config)
	if err != nil {
		return errors.Errorf("Failed to start Sarama producer: %v", err)
	}

	w.producer = producer

	return nil
}

func (w *walker) SetServerRouter(router *gin.Engine) {
	w.server.Handler = router
}

func (w *walker) Run() (e error) {
	time.Sleep(5 * time.Second)
	if e = w.walkerWillStart(); e != nil {
		log.Errorf("start walker failed, error: %s", e)
		return
	}

	time.Sleep(5 * time.Second)
	if e = w.walkerStarted(); e != nil {
		log.Errorf("start walker failed, error: %s", e)
		return
	}

	<-w.stopCh
	log.Debug("get stop signal from stopCh")
	return
}

func (w *walker) AppendJob(job *atypo.JobInfo) error {
	wd, e := database.Daos.Walker.GetWalker(w.name)
	if e != nil {
		return e
	}

	if len(wd.Jobs) > 100 {
		log.Warnf("walker is handling to many jobs,"+
			" another job is coming while processing %d now", len(wd.Jobs))
	}

	for _, n := range wd.Jobs {
		if n == job.Name {
			log.Warnf("job already exists")
			return nil
		}
	}

	if len(wd.Jobs) == 0 {
		wd.Jobs = []string{job.Name}
	} else {
		wd.Jobs = append(wd.Jobs, job.Name)
	}
	if e := database.Daos.Walker.UpdateWalker(wd); e != nil {
		return e
	}

	njob, e := database.Daos.Job.GetJobInfo(job.Name, job.UID)
	if e != nil {
		// we should never be here
		log.Warnf("fail to try to update job info, error: %v", e)
		return w.doneJob(job.Name)
	}

	njob.Status = atypo.RunningJobStatus
	// e = database.Daos.Job.Update

	if job.Params.FromFileList {
		return w.handleSocketJob(job)
	}
	return w.handleWalkJob(job)
}

func (w *walker) doneJob(jobName string) error {
	wd, e := database.Daos.Walker.GetWalker(w.name)
	if e != nil {
		return e
	}

	index := common.Index(wd.Jobs, jobName)

	if index == -1 {
		return errors.Errorf("%s not found in walker %s", jobName, w.name)
	}

	wd.Jobs = common.Filter(wd.Jobs, func(n string) bool {
		return n != jobName
	})

	if e := database.Daos.Walker.UpdateWalker(wd); e != nil {
		return e
	}

	return nil
}

func (w *walker) handleSocketJob(job *atypo.JobInfo) error {
	// TODO implement socket job
	return nil
}

func (w *walker) handleWalkJob(job *atypo.JobInfo) error {
	// TODO implement walk job
	go func() {
		w.listAlluxioPath(job, job.Params.AlluxioURI)
		log.Infof("walk job %s done", job.Name)
		w.doneJob(job.Name)
	}()
	return nil
}

func (w *walker) wcIncrease() {
	w.wcWg.Lock()
	w.walkCount++
	w.wcWg.Unlock()
}

func (w *walker) wcDecrease() {
	w.wcWg.Lock()
	w.walkCount--
	w.wcWg.Unlock()
}

func (w *walker) listAlluxioPath(job *atypo.JobInfo, path string) {
	w.wcWg.RLock()
	if w.walkCount < WALKER_CONCURRENT_LIMIT {
		w.wcWg.RUnlock()

		name := job.Name
		w.wcIncrease()

		result := &[]typo.AlluxioPath{}
		retryTime := 3
	retry:
		ctx := context.Background()
		e := w.rpcClient.Call(ctx,
			result, "GET", "http://"+w.primaryHost+"/api/v1/master/file/list_status?path="+path)

		if e != nil {
			log.Warnf("failed to list status of %s for job %s, error: e", path, job.Name, e)
		}

		if e != nil && retryTime > 0 {
			retryTime--
			goto retry
		}

		for index, item := range *result {
			if index%1000 == 0 {
				fmt.Printf("list %d for %s job\n", index, name)
			}
			if item.Folder {
				w.listAlluxioPath(job, item.Path)
				continue
			}

			if item.InAlluxioPercentage == 100 && job.Type == atypo.PreloadJobType {
				// pass
			} else if item.Persisted && job.Type == atypo.SaveJobType {
				// pass
			} else {
				msgData := &common.KafkaMessage{
					MsgType: common.AvioCMDMsg,
					AvioCMDData: common.AvioCMDData{
						JobType: string(job.Type),
						Path:    item.Path,
					},
				}
				msgStr, e := json.Marshal(msgData)
				if e != nil {
					// we should never be here
					log.Warnf("json stringify job msg failed: %v, error: %v, but we should never be here", msgStr, e)
				}
				msg := &kafka.ProducerMessage{
					Topic: w.conf.KAFKA.Topic,
					Value: kafka.StringEncoder(msgStr),
				}
				w.producer.SendMessage(msg)
			}
		}

		w.wcDecrease()
	} else {
		w.wcWg.RUnlock()
	}
}

func (w *walker) GracefullyShutdown() (e error) {
	log.Info("gracefully shutting down")
	if e = w.walkerWillStop(); e != nil {
		return
	}

	if e = w.walkerStopped(); e != nil {
		return
	}
	return nil
}

func (w *walker) walkerWillStart() error {
	log.Info("walker is starting")
	go (func(w *walker) {
		log.Warnf("server is %v", w.server)
		w.server.ListenAndServe()
	})(w)

	wd := &typo.Walker{
		Name:   common.Name(w.name),
		Status: typo.PreOnLine,
		Jobs:   make([]string, 0, 100),
	}
	if e := database.Daos.Walker.InsertWalker(wd); e != nil {
		w.server.Shutdown(context.Background())
		return e
	}

	return nil
}

func (w *walker) walkerStarted() error {
	log.Info("walker is started")

	wd := &typo.Walker{
		Name:   common.Name(w.name),
		Status: typo.OnLine,
	}

	return database.Daos.Walker.UpdateWalker(wd)
}

func (w *walker) walkerWillStop() error {
	log.Info("walker is shuting down")

	wd := &typo.Walker{
		Name:   common.Name(w.name),
		Status: typo.OffLine,
	}

	// TODO check whether this is the right place to shut down server
	go (func() {
		w.server.Shutdown(context.Background())
	})()

	return database.Daos.Walker.UpdateWalker(wd)
}

func (w *walker) walkerStopped() error {
	defer w.mgoSession.Close()

	for true {
		// wait util current jobs finished
		wd, e := database.Daos.Walker.GetWalker(w.name)
		if e != nil {
			return e
		}

		l := len(wd.Jobs)
		if l > 2 {
			time.Sleep(10 * time.Second)
		} else if l == 1 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	e := database.Daos.Walker.DeleteWalker(w.name)
	log.Info("about to send signal to stopCh")
	w.stopCh <- true
	return e
}
