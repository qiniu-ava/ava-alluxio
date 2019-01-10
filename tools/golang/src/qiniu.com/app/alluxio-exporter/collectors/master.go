package collectors

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"qiniu.com/app/alluxio-exporter/typo"
)

type MasterCollector struct {
	Host string

	TotalCapacity prometheus.Gauge

	TotalCapacityUsed prometheus.Gauge

	LostWorkers *prometheus.GaugeVec

	Workers *prometheus.GaugeVec

	MasterMetrics *prometheus.GaugeVec

	StartTmieMS prometheus.Gauge

	RunningWorker prometheus.Gauge

	HDDCapacity prometheus.Gauge

	HDDCapacityUsed prometheus.Gauge

	MemCapacity prometheus.Gauge

	MemCapacityUsed prometheus.Gauge

	SSDCapacity prometheus.Gauge

	SSDCapacityUsed prometheus.Gauge

	UFSCapacity prometheus.Gauge

	UFSCapacityUsed prometheus.Gauge
}

func NewMasterCollector(host string, group string) *MasterCollector {
	labels := make(prometheus.Labels)
	labels["group"] = group

	return &MasterCollector{
		Host: host,
		RunningWorker: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_live_worker_count",
			Help:        "Total number of the living worker",
			ConstLabels: labels,
		}),
		TotalCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_capacity_bytes_all",
			Help:        "Total capacity of the master",
			ConstLabels: labels,
		}),
		TotalCapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_capacity_bytes_used",
			Help:        "Used capacity of the master",
			ConstLabels: labels,
		}),
		MemCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_mem_capacity_bytes_all",
			Help:        "Total memory capacity of the master",
			ConstLabels: labels,
		}),
		MemCapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_mem_capacity_bytes_used",
			Help:        "Used memory capacity of the master",
			ConstLabels: labels,
		}),
		SSDCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_ssd_capacity_bytes_all",
			Help:        "Total SSD capacity of the master",
			ConstLabels: labels,
		}),
		SSDCapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_ssd_capacity_bytes_used",
			Help:        "Used SSD capacity of the master",
			ConstLabels: labels,
		}),
		HDDCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_hdd_capacity_bytes_all",
			Help:        "Total HDD capacity of the master",
			ConstLabels: labels,
		}),
		HDDCapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_hdd_capacity_bytes_used",
			Help:        "Used HDD capacity of the master",
			ConstLabels: labels,
		}),
		UFSCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_ufs_capacity_bytes_all",
			Help:        "Total UFS capacity of the master",
			ConstLabels: labels,
		}),
		UFSCapacityUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_ufs_capacity_bytes_used",
			Help:        "Used UFS capacity of the master",
			ConstLabels: labels,
		}),
		StartTmieMS: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   typo.AlluxioNamespace,
			Name:        "master_start_time",
			Help:        "Start time of the master",
			ConstLabels: labels,
		}),
		MasterMetrics: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.AlluxioNamespace,
				Name:        "master_metrics",
				Help:        "Metrics of the master",
				ConstLabels: labels,
			},
			[]string{"metric"},
		),
		Workers: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.AlluxioNamespace,
				Name:        "master_worker_conntected_state",
				Help:        "Status of the worker connected to the master",
				ConstLabels: labels,
			},
			[]string{"host", "state", "metric"},
		),
		LostWorkers: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.AlluxioNamespace,
				Name:        "master_worker_lost_state",
				Help:        "Status of the worker connected to the master",
				ConstLabels: labels,
			},
			[]string{"host", "state", "metric"},
		),
	}
}

func (m *MasterCollector) metricsList() []prometheus.Metric {
	return []prometheus.Metric{
		m.TotalCapacity,
		m.TotalCapacityUsed,
		m.MemCapacity,
		m.MemCapacityUsed,
		m.SSDCapacity,
		m.SSDCapacityUsed,
		m.StartTmieMS,
		m.HDDCapacity,
		m.HDDCapacityUsed,
		m.UFSCapacity,
		m.UFSCapacityUsed,
		m.RunningWorker,
	}
}

func (m *MasterCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		m.MasterMetrics,
		m.Workers,
		m.LostWorkers,
	}
}

func (m *MasterCollector) collect() error {
	scheme := "http://"
	host := m.Host
	path := "/api/v1/master/info/"
	method := "GET"
	URLStr := scheme + host + path

	res, e := HTTPRequest(URLStr, method, nil, nil)
	if e != nil {
		log.Println("[ERROR] HTTPRequest for master failed: ", e)
		return e
	}

	result := typo.MasterStat{}
	e = json.Unmarshal(res, &result)
	if e != nil {
		log.Println("[ERROR] json Unmarshal error:", e)
		return e
	}

	m.TotalCapacity.Set(result.Capacity.Total)
	m.TotalCapacityUsed.Set(result.Capacity.Used)
	m.MemCapacity.Set(result.CapacityAll.MEM.Total)
	m.MemCapacityUsed.Set(result.CapacityAll.MEM.Used)
	m.SSDCapacity.Set(result.CapacityAll.SSD.Total)
	m.SSDCapacityUsed.Set(result.CapacityAll.SSD.Used)
	m.UFSCapacity.Set(result.UFSCapacity.Total)
	m.UFSCapacityUsed.Set(result.UFSCapacity.Used)
	m.HDDCapacity.Set(result.CapacityAll.HDD.Total)
	m.HDDCapacityUsed.Set(result.CapacityAll.HDD.Used)
	m.StartTmieMS.Set(result.StartTimeMs)
	m.RunningWorker.Set(float64(len(result.Workers)))

	for k, v := range result.Metric {
		tmpName := strings.Split(k, ".")[1]
		m.MasterMetrics.WithLabelValues(tmpName).Set(v.(float64))
	}

	m.describeWorker(result.Workers)

	m.describeWorker(result.LostWorkers)

	return nil
}

func (m *MasterCollector) describeWorker(workers []typo.Worker) {
	for _, v := range workers {
		workerPort := strconv.Itoa(v.Address.DataPort)
		workerHost := v.Address.Host + ":" + workerPort
		m.Workers.WithLabelValues(workerHost, v.State, "capacityBytes").Set(v.Capacity)
		m.Workers.WithLabelValues(workerHost, v.State, "ID").Set(v.ID)
		m.Workers.WithLabelValues(workerHost, v.State, "lastConnectSec").Set(v.LastContactSec)
		m.Workers.WithLabelValues(workerHost, v.State, "startTimeMs").Set(v.StartTimeMs)
		m.Workers.WithLabelValues(workerHost, v.State, "usedBytes").Set(v.UsedBytes)
	}
}

func (m *MasterCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range m.metricsList() {
		ch <- metric.Desc()
	}
	for _, metric := range m.collectorList() {
		metric.Describe(ch)
	}
}

func (m *MasterCollector) Collect(ch chan<- prometheus.Metric) {
	if err := m.collect(); err != nil {
		log.Println("[ERROR] failed collecting cluster usage metrics:", err)
		return
	}

	for _, metric := range m.collectorList() {
		metric.Collect(ch)
	}

	for _, metric := range m.metricsList() {
		ch <- metric
	}
}
