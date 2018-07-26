package database

import (
	mgo "gopkg.in/mgo.v2"
	apiserverdb "qiniu.com/app/avio-apiserver/database"
	"qiniu.com/app/avio-walker/utils"
	log "qiniupkg.com/x/log.v7"
)

var db *mgo.Database

type daos struct {
	Job    apiserverdb.JobDao
	Walker WalkerDao
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

	j, e := apiserverdb.NewJobDao(db)
	if e != nil {
		log.Warn("init job collection failed")
		return e
	}
	Daos.Job = j

	w, e := NewWalkerDao(db)
	if e != nil {
		log.Warn("init walker collection failed")
		return e
	}
	Daos.Walker = w

	return
}
