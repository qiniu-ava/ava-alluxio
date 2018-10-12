package collectors

import (
	"encoding/json"
	"strconv"
	"testing"

	"qiniu.com/app/alluxio-exporter/typo"
)

func Test_WorkerInfo(t *testing.T) {
	res, err := HTTPRequest("http://127.0.0.1:30000/api/v1/worker/info/", "GET", nil, nil)
	if err != nil {
		t.Error(err)
	} else {
		result := typo.WorkerStat{}
		if e := json.Unmarshal(res, &result); e != nil {
			t.Error(err)
		}
		t.Logf("%v\n", result)
	}
}

func Test_WorkerMetric(t *testing.T) {

	hostIP := "127.0.0.1"

	res, err := HTTPRequest("http://"+hostIP+":30000/api/v1/worker/metrics/", "GET", nil, nil)
	if err != nil {
		t.Error(err)
	} else {
		result := make(map[string]int)
		if e := json.Unmarshal(res, &result); e != nil {
			t.Error(err)
		}
		for k, v := range result {
			t.Log(k + ":" + strconv.Itoa(v))
		}
	}
}
