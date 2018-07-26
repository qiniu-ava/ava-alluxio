package common

import (
	"regexp"

	"github.com/pkg/errors"
)

type Name string

var nameReg = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9-_]{3,10}")

func (n *Name) String() string {
	return string(*n)
}

func (n *Name) Validate() error {
	if !nameReg.MatchString(n.String()) {
		return errors.Errorf("invalid name")
	}
	return nil
}

type MsgType string

const (
	AvioCMDMsg MsgType = "avio_cmd_msg"
)

func (t *MsgType) Validate() error {
	var validTypes = []MsgType{AvioCMDMsg}
	for _, item := range validTypes {
		if string(*t) == string(item) {
			return nil
		}
	}

	return errors.Errorf("invalid kafka msg type: %v", t)
}

type AvioCMDData struct {
	JobType string `json:"jobType"`
	Path    string `json:"path"`
}

type KafkaMessage struct {
	MsgType     MsgType     `json:"msgType"`
	AvioCMDData AvioCMDData `json:"avioCMDData,omitempty"`
}
