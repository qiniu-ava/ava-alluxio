package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"qiniu.com/app/avio-apiserver/utils"
	"qiniu.com/app/common/database"
	"qiniu.com/app/common/typo"
	log "qiniupkg.com/x/log.v7"
)

type JobService struct {
	client http.Client
}

func NewJobService() *JobService {
	return &JobService{
		client: http.Client{
			Timeout: time.Duration(60 * time.Second),
		},
	}
}

func (j *JobService) CreateJob(c *gin.Context) {
	job := &typo.JobSpec{}
	if e := c.BindJSON(job); e != nil {
		log.Warnf("create job failed for invalid job spec, %v", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}

	log.Debugf("trying to insert job with data: %v", job)

	if e := job.Validate(); e != nil {
		log.Warnf("create job failed for invalid job spec, %v", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}

	info, e := database.Daos.Job.InsertJob(job)

	if e != nil {
		log.Warnf("insert job error: %v", e)
		c.JSON(http.StatusNotImplemented, gin.H{
			"err": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": job.Name,
		"item": info,
	})
}

func (j *JobService) StartJob(c *gin.Context) {
	jobName := c.Param("name")
	u := c.GetHeader("X-UID")
	uid, e := strconv.ParseInt(u, 10, 64)
	if e != nil {
		log.Warnf("uid string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}

	job, e := database.Daos.Job.GetJobInfo(jobName, uid)
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": fmt.Sprintf("获取任务信息失败"),
		})
		return
	}

	if job.Status != typo.CreatedJobStatus && job.Status != typo.FailedJobStatus {
		c.JSON(http.StatusForbidden, gin.H{
			"err": fmt.Sprintf("当前任务状态不允许再次启动"),
		})
		return
	}

	w, e := database.Daos.Walker.GetWalkers()
	if e != nil || len(w) == 0 {
		log.Warnf("no active walker available now")
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": fmt.Sprintf("后台服务器暂时不可用，请稍后再试"),
		})
	}

	walkerName := utils.PickWalker(w)

	req := http.Request{
		URL: &url.URL{
			Host:   walkerName,
			Scheme: "http",
			Path:   "/jobs/" + jobName + "/start",
		},
		Method: "POST",
		Header: http.Header{
			"X-UID": []string{fmt.Sprintf("%d", uid)},
		},
	}

	timesToRetry := 3
retry:
	res, e := j.client.Do(&req)
	if e != nil {
		if timesToRetry > 0 {
			timesToRetry--
			log.Errorf("start job %s on walk server %s failed after %d retries, error: %v", jobName, walkerName, 3-timesToRetry, e)
			goto retry
		} else {
			log.Errorf("start job %s on walk server %s failed after %d retries, error: %v", jobName, walkerName, 3-timesToRetry, e)
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": e,
			})
			return
		}
	}

	if res.StatusCode-res.StatusCode%100 != 200 {
		log.Errorf("start job %s on walk server %s failed, status code: %d", jobName, walkerName, res.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": fmt.Sprintf("服务器错误，错误码%d", res.StatusCode),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (j *JobService) ListJobs(c *gin.Context) {
	limit, e := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if e != nil {
		log.Warnf("limit string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}

	skip, e := strconv.Atoi(c.DefaultQuery("skip", "0"))
	if e != nil {
		log.Warnf("skip string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}
	uid, e := strconv.ParseInt(c.Query("uid"), 10, 64)
	if e != nil {
		log.Warnf("uid string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}

	query := &typo.ListJobQuery{
		Limit: limit,
		Skip:  skip,
		UID:   uid,
	}

	result := &typo.ListJobResult{}

	if e := database.Daos.Job.ListJob(result, query); e != nil {
		log.Warnf("list jobs error: %v", e)
		c.JSON(http.StatusNotImplemented, gin.H{
			"err": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query": gin.H{
			"limit": query.Limit,
			"skip":  query.Skip,
			"uid":   query.UID,
		},
		"total": result.Total,
		"items": result.Items,
	})

}

func (j *JobService) GetJobInfo(c *gin.Context) {
	jobName := c.Param("name")
	headUID := c.GetHeader("X-UID")
	if headUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "no uid in header",
		})
		return
	}
	uid, e := strconv.ParseInt(headUID, 10, 64)
	if e != nil {
		log.Warnf("uid string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}
	if e != nil {
		log.Warnf("uid string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}
	info, e := database.Daos.Job.GetJobInfo(jobName, uid)
	if e != nil {
		log.Warnf("get job info error: %v", e)
		c.JSON(http.StatusNotImplemented, gin.H{
			"err": e,
		})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (j *JobService) UpdateJobInfo(c *gin.Context) {
	jobName := c.Param("name")

	headUID := c.GetHeader("X-UID")
	if headUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "no uid in header",
		})
		return
	}
	uid, e := strconv.ParseInt(headUID, 10, 64)
	if e != nil {
		log.Warnf("uid string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}

	jobInfo := &typo.JobInfo{}
	if e := c.BindJSON(jobInfo); e != nil {
		log.Warnf("update job failed for invalid job info, %v", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}

	log.Debugf("trying to update job with data: %v", jobInfo)

	if e := jobInfo.Validate(); e != nil {
		log.Warnf("update job failed for invalid job info, %v", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}

	updatedInfo, e := database.Daos.Job.UpdateJobInfo(jobName, uid, jobInfo)

	if e != nil {
		log.Warnf("update job error: %v", e)
		c.JSON(http.StatusNotImplemented, gin.H{
			"err": e,
		})
		return
	}

	c.JSON(http.StatusOK, updatedInfo)
}

func (j *JobService) DeleteJobInfo(c *gin.Context) {
	jobName := c.Param("name")
	headUID := c.GetHeader("X-UID")
	if headUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "no uid in header",
		})
		return
	}
	uid, e := strconv.ParseInt(headUID, 10, 64)
	if e != nil {
		log.Warnf("uid string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}
	if e := database.Daos.Job.DeleteJobInfo(jobName, uid); e != nil {
		log.Warnf("delete job error: %v", e)
		c.JSON(http.StatusNotImplemented, gin.H{
			"err": e,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
