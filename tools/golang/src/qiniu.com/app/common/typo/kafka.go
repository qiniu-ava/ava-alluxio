package typo

import (
	"github.com/pkg/errors"
)

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
	JobType JobType `json:"jobType"`
	Path    string  `json:"path"`
}

type KafkaMessage struct {
	MsgType     MsgType     `json:"msgType"`
	AvioCMDData AvioCMDData `json:"avioCMDData,omitempty"`
}
