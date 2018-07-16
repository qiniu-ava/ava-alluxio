package typo

import (
	"time"

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
	AlluxioUri   string `json: "alluxioUri"`
	Depth        int    `json: "depth"`
	FromFileList bool   `json: "fromFileList"`
}

type JobSpec struct {
	Name   string    `json:"name"`
	Uid    int       `json: "uid"`
	Type   JobType   `json: "type"`
	Params JobParams `json: "params"`
}

type JobInfo struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	JobSpec    `bson:",inline"`
	Status     JobStatus  `json: "status"`
	CreateTime *time.Time `json:"createTime,omitempty" bson:"createTime,omitempty"`
	UpdateTime *time.Time `json:"updateTime,omitempty" bson:"updateTime,omitempty"`
	Message    string     `json:"message,omitempty" bson:"message,omitempty"`
}

func (t *JobType) Validate() error {
	// TODO implement me
	return nil
}

func (p *JobParams) Validate() error {
	// TODO implement me
	return nil
}

func (s *JobSpec) Validate() error {
	// TODO implement me
	return nil
}
