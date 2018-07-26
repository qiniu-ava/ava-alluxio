package services

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"qiniu.com/app/avio-walker/database"
	"qiniu.com/app/avio-walker/walker"
	log "qiniupkg.com/x/log.v7"
)

type WalkJobService interface {
	StartWalkJob(c *gin.Context)
}

type walkJobService struct {
	walker walker.Walker
}

func NewWalkJobService(w walker.Walker) WalkJobService {
	return &walkJobService{
		walker: w,
	}
}

func (w *walkJobService) StartWalkJob(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		log.Warnf("bad request, empty name")
		return
	}

	u := c.GetHeader("X-UID")
	if u == "" {
		log.Warnf("bad request, header X-UID is empty")
		return
	}

	uid, e := strconv.Atoi(u)

	if e != nil {
		log.Warnf("bad request, invalid header X-UID")
		return
	}

	job, e := database.Daos.Job.GetJobInfo(name, uid)

	if e != nil {
		log.Warnf("bad request, job not found")
		return
	}

	if e := w.walker.AppendJob(job); e != nil {
		log.Warnf("bad request, append job failed")
		return
	}

	return
}
