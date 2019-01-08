package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"qiniu.com/app/alluxio-exporter/collectors"
)

type AlluxioExporter struct {
	mu         sync.Mutex
	collectors []prometheus.Collector
}

var _ prometheus.Collector = &AlluxioExporter{}

func NewAlluxioExporter() *AlluxioExporter {

	return &AlluxioExporter{}
}

func (a *AlluxioExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, cc := range a.collectors {
		cc.Describe(ch)
	}
}

func (a *AlluxioExporter) Collect(ch chan<- prometheus.Metric) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, cc := range a.collectors {
		cc.Collect(ch)
	}
}

func main() {
	var (
		addr           = flag.String("telemetry.addr", ":9996", "host:port for alluxio exporter")
		metricsPath    = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")
		exporterConfig = flag.String("exporter.config", "/conf/alluxio-exporter.yml", "Path to alluxio exporter config.")
	)
	flag.Parse()
	log.Printf("parse the file ", *exporterConfig)
	exporter := NewAlluxioExporter()

	if fileExists(*exporterConfig) {
		cfg, err := ParseConfig(*exporterConfig)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		for _, alluxio := range cfg.Alluxio {
			group := alluxio.Group
			log.Printf("Group name is :", group)
			exporter.collectors = append(exporter.collectors,
				collectors.NewMasterCollector(alluxio.MasterHost, group))
			workerList, err := GetWorkerInfoList(alluxio.MasterHost)
			if err != nil {
				log.Fatalf("Get worker info Error: %v", err)
			}
			for _, worker := range workerList {
				instance := worker.Address.Host + ":" + strconv.Itoa(worker.Address.WebPort)
				role := worker.Address.Role
				log.Printf("[INFO]Worker instance : ", instance)
				exporter.collectors = append(exporter.collectors,
					collectors.NewWorkerCollector(instance, group, role, worker.Address.Host))
			}
		}
	} else {
		log.Printf("file do not exist !")
		return
	}
	prometheus.MustRegister(exporter)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Alluxio Exporter</title></head>
			<body>
			<h1>Alluxio Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Starting alluxio exporter on %q", *addr)
	ln, err := net.Listen("tcp", *addr)
	if err == nil {
		err := http.Serve(ln.(*net.TCPListener), nil)
		if err != nil {
			log.Fatalf("unable to serve requests: %s", err)
		}
	}
	if err != nil {
		log.Fatalf("unable to create listener: %s", err)
	}
}
