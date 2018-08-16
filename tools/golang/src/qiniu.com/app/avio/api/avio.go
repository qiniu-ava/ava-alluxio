package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"qiniu.com/app/avio/constants"
	"qiniu.com/app/common/typo"
)

func NewJob(client http.Client, job typo.JobSpec) error {
	jsonBytes, _ := json.Marshal(&job)
	url := "http://" + constants.AVIO_SERVICE_HOST + "/jobs"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("X-UID", fmt.Sprintf("%d", job.UID))
	timeToRetry := 3
retry:
	res, e := client.Do(req)
	if e != nil && timeToRetry > 0 {
		if timeToRetry > 0 {
			goto retry
		} else {
			if res.StatusCode%100 == 5 {
				return fmt.Errorf("网络请求失败，服务器端故障。")
			}
			return fmt.Errorf("网络请求失败，请确定你的网络状况良好。")
		}
	}
	return nil
}

func StartJob(client http.Client, jobName string, uid int64) error {
	url := "http://" + constants.AVIO_SERVICE_HOST + "/jobs/" + jobName + "/start"
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("X-UID", fmt.Sprintf("%d", uid))
	timeToRetry := 3
retry:
	res, e := client.Do(req)
	if e != nil && timeToRetry > 0 {
		if timeToRetry > 0 {
			goto retry
		} else {
			if res.StatusCode%100 == 5 {
				return fmt.Errorf("网络请求失败，服务器端故障。")
			}
			return fmt.Errorf("网络请求失败，请确定你的网络状况良好。")
		}
	}
	return nil
}

func GetJobs(client http.Client, uid int64, limit int, skip int) error {
	return nil
}

func GetJob(jobName string, uid int64) error {
	return nil
}
