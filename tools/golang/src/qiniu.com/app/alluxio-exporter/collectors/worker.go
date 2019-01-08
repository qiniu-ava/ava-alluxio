package collectors

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"qiniu.com/app/alluxio-exporter/typo"
)

type WorkerCollector struct {
	instance string

	host string

	role string

	TotalCapacity prometheus.Gauge

	TotalCapacityUsed prometheus.Gauge

	StartTmieMS prometheus.Gauge

	MemCapacity prometheus.Gauge

	MemCapacityUsed prometheus.Gauge

	SSDCapacity prometheus.Gauge

	SSDCapacityUsed prometheus.Gauge

	WorkerMetrics *prometheus.GaugeVec

	Roles *prometheus.GaugeVec
}

func NewWorkerCollector(workerAddress string, group string, role string, host string) *WorkerCollector {
	labels := make(prometheus.Labels)
	labels["instance"] = workerAddress
	labels["group"] = group

	return &WorkerCollector{
		instance: workerAddress,
		role:     role,
		host:     host,
		TotalCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "worker_capacity_bytes_all",
			Help:        "Total capacity of the worker",
			ConstLabels: labels,
		}),
		TotalCapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "worker_capacity_bytes_used",
			Help:        "Used capacity of the worker",
			ConstLabels: labels,
		}),
		MemCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "worker_mem_capacity_bytes_all",
			Help:        "Total memory capacity of the worker",
			ConstLabels: labels,
		}),
		MemCapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "worker_mem_capacity_bytes_used",
			Help:        "Used memory capacity of the worker",
			ConstLabels: labels,
		}),
		SSDCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "worker_ssd_capacity_bytes_all",
			Help:        "Total SSD capacity of the worker",
			ConstLabels: labels,
		}),
		SSDCapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "worker_ssd_capacity_bytes_used",
			Help:        "Used SSD capacity of the worker",
			ConstLabels: labels,
		}),
		StartTmieMS: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "worker_start_time",
			Help:        "Start time of the worker",
			ConstLabels: labels,
		}),
		WorkerMetrics: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.AlluxioNamespace,
				Name:        "worker_metrics",
				Help:        "Metrics of the worker",
				ConstLabels: labels,
			},
			[]string{"metric"},
		),
		Roles: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.AlluxioNamespace,
				Name:        "worker_roles",
				Help:        "Roles of the worker",
				ConstLabels: labels,
			},
			[]string{"host", "role"},
		),
	}
}

func (w *WorkerCollector) metricsList() []prometheus.Metric {
	return []prometheus.Metric{
		w.TotalCapacity,
		w.TotalCapacityUsed,
		w.MemCapacity,
		w.MemCapacityUsed,
		w.SSDCapacity,
		w.SSDCapacityUsed,
		w.StartTmieMS,
	}
}

func (w *WorkerCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		w.WorkerMetrics,
		w.Roles,
	}
}

func (w *WorkerCollector) collect() error {
	scheme := "http://"
	host := w.instance
	path := "/api/v1/worker/info"
	method := "GET"
	URLStr := scheme + host + path

	res, e := HTTPRequest(URLStr, method, nil, nil)
	if e != nil {
		log.Println("[ERROR] HTTPRequest for worker:"+w.instance+" failed: ", e)
		return e
	}

	result := typo.WorkerStat{}
	e = json.Unmarshal(res, &result)
	if e != nil {
		log.Println("[ERROR] json Unmarshal error:", e)
		return e
	}

	w.TotalCapacity.Set(result.Capacity.Total)
	w.TotalCapacityUsed.Set(result.Capacity.Used)
	w.MemCapacity.Set(result.CapacityAll.MEM.Total)
	w.MemCapacityUsed.Set(result.CapacityAll.MEM.Used)
	w.SSDCapacity.Set(result.CapacityAll.SSD.Total)
	w.SSDCapacityUsed.Set(result.CapacityAll.SSD.Used)
	w.StartTmieMS.Set(result.StartTimeMs)

	for k, v := range result.Metric {
		tmpName := strings.Split(k, ".")[2]
		w.WorkerMetrics.WithLabelValues(tmpName).Set(v.(float64))
	}
	w.describeRoles()
	return nil
}

func (w *WorkerCollector) describeRoles() {
	w.Roles.WithLabelValues(w.host, w.role)
}

func (w *WorkerCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range w.metricsList() {
		ch <- metric.Desc()
	}
	for _, metric := range w.collectorList() {
		metric.Describe(ch)
	}
}

func (w *WorkerCollector) Collect(ch chan<- prometheus.Metric) {
	if err := w.collect(); err != nil {
		log.Println("[ERROR] failed collecting cluster usage metrics:", err)
		return
	}

	for _, metric := range w.collectorList() {
		metric.Collect(ch)
	}

	for _, metric := range w.metricsList() {
		ch <- metric
	}
}
