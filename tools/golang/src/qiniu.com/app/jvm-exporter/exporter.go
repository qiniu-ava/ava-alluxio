package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"qiniu.com/app/jvm-exporter/collectors"
)

type JVMExporter struct {
	mu         sync.Mutex
	collectors []prometheus.Collector
}

var _ prometheus.Collector = &JVMExporter{}

func NewJVMExporter() *JVMExporter {
	return &JVMExporter{}
}

func (j *JVMExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, cc := range j.collectors {
		cc.Describe(ch)
	}
}

func (j *JVMExporter) Collect(ch chan<- prometheus.Metric) {
	j.mu.Lock()
	defer j.mu.Unlock()

	for _, cc := range j.collectors {
		cc.Collect(ch)
	}
}

func main() {
	var (
		addr        = flag.String("telemetry.addr", ":9998", "host:port for alluxio exporter")
		metricsPath = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")
		shellPath   = flag.String("gcMonitor.Path", "/tools/gcMetric.sh", "Path to JVM GC shell path.")
	)
	flag.Parse()
	ipAddr := os.Getenv("HOSTIP")
	if ipAddr == "" {
		err := errors.Errorf("Host IP should be supplied")
		log.Fatalf("Error: %v", err)
	}
	exporter := NewJVMExporter()
	stat, err := os.Stat(*shellPath)
	if !os.IsNotExist(err) && !stat.IsDir() {
		exporter.collectors = append(exporter.collectors,
			collectors.NewGcCollector(ipAddr, *shellPath))
	} else {
		err := errors.Errorf("GC shell path should be supplied and existed")
		log.Fatalf("Error: %v", err)
	}
	prometheus.MustRegister(exporter)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>JVM Exporter</title></head>
			<body>
			<h1>Alluxio Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Starting JVM exporter on %q", *addr)
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
