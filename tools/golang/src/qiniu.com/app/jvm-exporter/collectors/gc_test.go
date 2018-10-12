package collectors

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func Test_GCMetric(t *testing.T) {
	cmd := exec.Command("/bin/bash", "-c",
		"/Users/qnxr/Documents/alluxio/ava-alluxio/ava-alluxio/tools/golang/src/qiniu.com/app/jvm-exporter/gcMetric.sh")
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		t.Error(err)
	}
	ba := bytes.Split(out.Bytes(), []byte("\n"))
	m := make(map[string]string)
	for _, v := range ba {
		if len(v) < 1 {
			continue
		}
		if err = json.Unmarshal(v, &m); err != nil {
			t.Error(err)
		}
	}
	t.Logf("map is %v", m)
	for k, v := range m {
		if k == "alluxio-master" || k == "alluxio-worker" {
			vSlice := strings.Split(v, " ")
			var numbers []float64
			for _, elem := range vSlice {
				if elem == "" {
					continue
				}
				i, err := strconv.ParseFloat(elem, 64)
				if err != nil {
					t.Error(err)
				}
				numbers = append(numbers, i)
			}
			if len(numbers) < 5 {
				err = errors.Errorf("GC result length wrong")
				t.Error(err)
			}
		} else {
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				t.Error(err)
			}
			desList := strings.Split(k, " ")
			if len(desList) != 2 {
				err = errors.Errorf("GC Thread key length wrong")
				t.Error(err)
			}
			t.Logf("key is %s %s, value is %f", desList[0], desList[1], value)
		}
	}
}
