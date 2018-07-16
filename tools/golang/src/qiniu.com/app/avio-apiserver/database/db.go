package database

import (
	mgo "gopkg.in/mgo.v2"
	"qiniu.com/app/avio-apiserver/utils"
	log "qiniupkg.com/x/log.v7"
)

var db *mgo.Database

type daos struct {
	Job JobDao
}

var Daos *daos

func NewMongoSession(conf *utils.Config) (session *mgo.Session, e error) {
	session, e = mgo.Dial(conf.DB.Host)
	if e != nil {
		return
	}

	session.SetPoolLimit(conf.DB.SessionPoolSize)

	db = session.DB(conf.DB.DBName)
	return
}

func Init() (e error) {
	Daos = &daos{}
	j, e := NewJobDao()
	Daos.Job = j
	if e != nil {
		log.Warn("init job collection failed")
		return
	}

	return
}
