package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"qiniu.com/app/avio-apiserver/database"
	"qiniu.com/app/avio-apiserver/typo"
	log "qiniupkg.com/x/log.v7"
)

type JobService struct {
	dao database.JobDao
}

func NewJobService() *JobService {
	return &JobService{
		dao: database.Daos.Job,
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

	info, e := j.dao.InsertJob(job)

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
	uid, e := strconv.Atoi(c.Query("uid"))
	if e != nil {
		log.Warnf("uid string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}

	query := &database.ListJobQuery{limit, skip, uid}

	result := &database.ListJobResult{}

	if e := j.dao.ListJob(result, query); e != nil {
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
		"item":  result.Items,
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
	uid, e := strconv.Atoi(headUID)
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
	info, e := j.dao.GetJobInfo(jobName, uid)
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
	uid, e := strconv.Atoi(headUID)
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

	updatedInfo, e := j.dao.UpdateJobInfo(jobName, uid, jobInfo)

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
	uid, e := strconv.Atoi(headUID)
	if e != nil {
		log.Warnf("uid string2int error : ", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"err": e,
		})
		return
	}
	if e := j.dao.DeleteJobInfo(jobName, uid); e != nil {
		log.Warnf("delete job error: %v", e)
		c.JSON(http.StatusNotImplemented, gin.H{
			"err": e,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
