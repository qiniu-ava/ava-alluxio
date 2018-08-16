package typo

import (
	"regexp"
	"time"

	"github.com/pkg/errors"
)

type WalkerStatus int

const (
	PreOnLine WalkerStatus = 0
	OnLine    WalkerStatus = 1
	OffLine   WalkerStatus = 2
)

func (w *WalkerStatus) validate() error {
	if *w != PreOnLine && *w != OnLine && *w != OffLine {
		return errors.Errorf("invalide walker status")
	}
	return nil
}

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

type Walker struct {
	Name       Name         `json:"name" bson:"name"`
	Status     WalkerStatus `json:"status" bson:"status"`
	Jobs       []string     `json:"jobs" bson:"jobs"`
	CreateTime time.Time    `json:"createTime" bson:"createTime"`
}

func (w *Walker) Validate() error {
	if e := w.Status.validate(); e != nil {
		return e
	}

	if e := w.Name.Validate(); e != nil {
		return e
	}

	if w.Jobs == nil {
		return errors.Errorf("jobs in walker should not be nil")
	}

	return nil
}
