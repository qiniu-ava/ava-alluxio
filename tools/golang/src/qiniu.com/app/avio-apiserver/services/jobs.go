package services

import (
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
		log.Warn("create job failed for invalid job spec, %v", e)
		return
	}

	log.Debugf("trying to insert job with data: %v", job)

	if e := job.Validate(); e != nil {
		log.Warn("create job failed for invalid job spec, %v", e)
		return
	}

	info, e := j.dao.InsertJob(job)
	log.Warnf("insert job: %v", info)

	if e != nil {
		return
	}
}

func (j *JobService) ListJobs(c *gin.Context) {
}

func (j *JobService) GetJobInfo(c *gin.Context) {

}

func (j *JobService) UpdateJobInfo(c *gin.Context) {

}

func (j *JobService) DeleteJobInfo(c *gin.Context) {

}
