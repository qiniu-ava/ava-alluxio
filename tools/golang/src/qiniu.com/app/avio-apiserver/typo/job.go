package typo

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type JobType int

const (
	ListJobType    JobType = 0
	PreloadJobType JobType = 1
	SaveJobType    JobType = 2
	StatJobType    JobType = 3
)

type JobStatus string

const (
	CreatedJobStatus JobStatus = "created"
	RunningJobStatus JobStatus = "running"
	SuccessJobStatus JobStatus = "success"
	FailedJobStatus  JobStatus = "failed"
)

type JobParams struct {
	AlluxioURI   string `json:"alluxioUri" bson:"alluxioUri"`
	Depth        int    `json:"depth" bson:"depth"`
	FromFileList bool   `json:"fromFileList" bson:"fromFileList"`
}

type JobSpec struct {
	Name   string    `json:"name" bson:"name"`
	UID    int       `json:"uid" bson:"uid"`
	Type   JobType   `json:"type" bson:"type"`
	Params JobParams `json:"params" bson:"params"`
}

type JobInfo struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	JobSpec    `bson:",inline" bson:"jobSpec"`
	Status     JobStatus  `json:"status" bson:"status"`
	CreateTime *time.Time `json:"createTime,omitempty" bson:"createTime,omitempty"`
	UpdateTime *time.Time `json:"updateTime,omitempty" bson:"updateTime,omitempty"`
	FinishTime *time.Time `json:"finishTime,omitempty" bson:"finishTime,omitempty"`
	Message    string     `json:"message,omitempty" bson:"message,omitempty"`
}

func (t *JobType) Validate() error {
	if *t != ListJobType && *t != PreloadJobType && *t != SaveJobType && *t != StatJobType {
		return errors.Errorf("illegal jobType")
	}
	return nil
}

func (p *JobParams) Validate() error {
	if _, e := url.Parse(p.AlluxioURI); e != nil {
		return e
	}
	return nil
}

func (s *JobSpec) Validate() error {
	if e := s.Type.Validate(); e != nil {
		return e
	}
	if e := s.Params.Validate(); e != nil {
		return e
	}
	return nil
}

func (i *JobInfo) Validate() error {
	if e := i.JobSpec.Validate(); e != nil {
		return e
	}
	if i.Status != CreatedJobStatus && i.Status != RunningJobStatus &&
		i.Status != SuccessJobStatus && i.Status != FailedJobStatus {
		return errors.Errorf("illegal jobType")
	}
	return nil
}
