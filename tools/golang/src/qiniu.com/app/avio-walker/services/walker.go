package services

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"qiniu.com/app/avio-walker/walker"
	"qiniu.com/app/common/database"
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

	uid, e := strconv.ParseInt(u, 10, 64)

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
