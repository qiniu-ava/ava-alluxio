package database

import (
	"time"

	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"qiniu.com/app/common/typo"
	log "qiniupkg.com/x/log.v7"
)

const jobCollectionName string = "job"

type JobDao interface {
	InsertJob(spec *typo.JobSpec) (info *typo.JobInfo, e error)
	ListJob(result *typo.ListJobResult, query *typo.ListJobQuery) error
	GetJobInfo(name string, uid int64) (info *typo.JobInfo, e error)
	UpdateJobInfo(name string, uid int64, info *typo.JobInfo) (infoUpdated *typo.JobInfo, e error)
	DeleteJobInfo(name string, uid int64) error
}

type jobDao struct {
	collection *mgo.Collection
}

func NewJobDao(db *mgo.Database) (d *jobDao, e error) {
	c := db.C(jobCollectionName)
	if e := c.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true, Name: "job_name"}); e != nil {
		return nil, e
	}
	d = &jobDao{
		collection: c,
	}
	if e := d.collection.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true, Name: "job_name"}); e != nil {
		return nil, e
	}
	return d, nil
}

func (j *jobDao) InsertJob(spec *typo.JobSpec) (info *typo.JobInfo, e error) {
	cTime := time.Now()
	bTime := time.Unix(0, 0)

	info = &typo.JobInfo{
		Status:     typo.CreatedJobStatus,
		CreateTime: &cTime,
		UpdateTime: &cTime,
		FinishTime: &bTime,
		Message:    "",
	}

	info.JobSpec = *spec

	if e := j.collection.Insert(info); e != nil {
		return nil, e
	}
	info, e = j.GetJobInfo(spec.Name, spec.UID)
	if e != nil {
		return
	}
	return
}

func (j *jobDao) ListJob(result *typo.ListJobResult, query *typo.ListJobQuery) error {
	q := bson.M{"uid": query.UID}
	r := j.collection.Find(q)
	c, e := r.Count()
	if e != nil {
		return e
	} else if c == 0 {
		return errors.Errorf("no job with uid %s in mongodb", query.UID)
	} else if query.Skip > c {
		return errors.Errorf("skip too many, only %s jobs in all", c)
	}
	result.Query = *query
	result.Total = c
	item := typo.JobInfo{}
	iter := r.Skip(query.Skip).Limit(query.Limit).Iter()
	for iter.Next(&item) {
		result.Items = append(result.Items, item)
	}
	return nil
}

func (j *jobDao) GetJobInfo(name string, uid int64) (info *typo.JobInfo, e error) {
	q := bson.M{"name": name, "uid": uid}
	r := j.collection.Find(q)
	if c, e := r.Count(); e != nil || c != 1 {
		if c == 0 {
			e = errors.Errorf("no job named %s in mongodb", name)
		} else if c > 1 {
			e = errors.Errorf("more then one job named %s in mongodb, that's terrible", name)
		}
		return nil, e
	}
	e = r.One(&info)
	log.Debugf("get job info: %v", info)
	return
}

func (j *jobDao) UpdateJobInfo(name string, uid int64, info *typo.JobInfo) (infoUpdated *typo.JobInfo, e error) {
	cTime := time.Now()
	switch info.Status {
	case typo.RunningJobStatus:
		info.UpdateTime = &cTime
	case typo.SuccessJobStatus, typo.FailedJobStatus:
		info.UpdateTime = &cTime
		info.FinishTime = &cTime
	}
	q := bson.M{"name": name, "uid": uid}
	u := bson.M{"$set": info}
	if e := j.collection.Update(q, u); e != nil {
		return nil, e
	}
	infoUpdated, e = j.GetJobInfo(name, uid)
	if e != nil {
		return nil, e
	}
	return
}

func (j *jobDao) DeleteJobInfo(name string, uid int64) error {
	if _, e := j.GetJobInfo(name, uid); e != nil {
		return e
	}
	q := bson.M{"name": name, "uid": uid}
	if _, e := j.collection.RemoveAll(q); e != nil {
		return e
	}
	return nil
}
