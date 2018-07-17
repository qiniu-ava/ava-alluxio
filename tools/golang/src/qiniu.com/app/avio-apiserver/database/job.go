package database

import (
	"gopkg.in/mgo.v2"
	"qiniu.com/app/avio-apiserver/typo"
	log "qiniupkg.com/x/log.v7"
)

const jobCollectionName string = "job"

type ListJobQuery struct {
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
	UID   int `json:"uid"`
}

type ListJobResult struct {
	Query ListJobQuery   `json:"query"`
	Items []typo.JobInfo `json:"items"`
	Total int            `json:"total"`
}

type JobDao interface {
	InsertJob(spec *typo.JobSpec) (info *typo.JobInfo, e error)
	ListJob(result *ListJobResult, query *ListJobQuery) error
}

type jobDao struct {
	collection *mgo.Collection
}

func NewJobDao() (d *jobDao, e error) {
	d = &jobDao{
		collection: db.C(jobCollectionName),
	}
	return nil, nil
}

func (j *jobDao) InsertJob(spec *typo.JobSpec) (info *typo.JobInfo, e error) {
	// TODO implement me
	log.Debugf("debug: insert job, implement me")
	return
}

func (j *jobDao) ListJob(result *ListJobResult, query *ListJobQuery) error {
	return nil
}
