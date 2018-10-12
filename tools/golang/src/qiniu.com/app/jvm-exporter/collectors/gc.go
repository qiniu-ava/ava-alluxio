package collectors

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"qiniu.com/app/jvm-exporter/typo"
)

type GcCollector struct {
	instance string

	shellPath string

	YGNumber *prometheus.GaugeVec

	YGTime *prometheus.GaugeVec

	FGNumber *prometheus.GaugeVec

	FGTime *prometheus.GaugeVec

	GCTime *prometheus.GaugeVec

	GCThread *prometheus.GaugeVec
}

func NewGcCollector(host, sp string) *GcCollector {
	labels := make(prometheus.Labels)
	labels["endpoint"] = "jvm-export"
	labels["instance"] = host

	return &GcCollector{
		instance:  host,
		shellPath: sp,
		YGNumber: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.JvmNamespace,
				Name:        "young_gc_times",
				Help:        "Young GC times of alluxio components",
				ConstLabels: labels,
			},
			[]string{"name"},
		),
		YGTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.JvmNamespace,
				Name:        "young_gc_duration",
				Help:        "Young GC duration of alluxio components",
				ConstLabels: labels,
			},
			[]string{"name"},
		),
		FGNumber: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.JvmNamespace,
				Name:        "full_gc_times",
				Help:        "Full GC times of alluxio components",
				ConstLabels: labels,
			},
			[]string{"name"},
		),
		FGTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.JvmNamespace,
				Name:        "full_gc_duration",
				Help:        "Full GC duration of alluxio components",
				ConstLabels: labels,
			},
			[]string{"name"},
		),
		GCTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.JvmNamespace,
				Name:        "all_gc_duration",
				Help:        "All GC duration of alluxio components",
				ConstLabels: labels,
			},
			[]string{"name"},
		),
		GCThread: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   typo.JvmNamespace,
				Name:        "all_gc_thread_cpu_seconds",
				Help:        "All GC thread cpu used of alluxio components",
				ConstLabels: labels,
			},
			[]string{"name", "thread"},
		),
	}
}

func (g *GcCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		g.YGNumber,
		g.YGTime,
		g.FGNumber,
		g.FGTime,
		g.GCTime,
		g.GCThread,
	}
}

func (g *GcCollector) collect() error {
	cmd := exec.Command("/bin/sh", "-c", g.shellPath)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	ba := bytes.Split(out.Bytes(), []byte("\n"))
	m := make(map[string]string)
	for _, v := range ba {
		if len(v) < 1 {
			continue
		}
		if err = json.Unmarshal(v, &m); err != nil {
			return err
		}
	}
	for k, v := range m {
		if k == "alluxio-master" || k == "alluxio-worker" || k == "alluxio-worker-writer" {
			vSlice := strings.Split(v, " ")
			var numbers []float64
			for _, elem := range vSlice {
				if elem == "" {
					continue
				}
				i, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					return err
				}
				numbers = append(numbers, i)
			}
			if len(numbers) < 5 {
				err = errors.Errorf("GC result length wrong")
				return err
			}
			g.YGNumber.WithLabelValues(k).Set(numbers[0])
			g.YGTime.WithLabelValues(k).Set(numbers[1])
			g.FGNumber.WithLabelValues(k).Set(numbers[2])
			g.FGTime.WithLabelValues(k).Set(numbers[3])
			g.GCTime.WithLabelValues(k).Set(numbers[4])
		} else {
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			desList := strings.Split(k, " ")
			if len(desList) != 2 {
				err = errors.Errorf("GC Thread key length wrong")
				return err
			}
			g.GCThread.WithLabelValues(desList[0], desList[1]).Set(value)
		}
	}
	return nil
}

func (g *GcCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range g.collectorList() {
		metric.Describe(ch)
	}
}

func (g *GcCollector) Collect(ch chan<- prometheus.Metric) {
	if err := g.collect(); err != nil {
		log.Println("[ERROR] failed collecting GC metrics:", err)
		return
	}

	for _, metric := range g.collectorList() {
		metric.Collect(ch)
	}
}
