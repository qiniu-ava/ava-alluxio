package database

import (
	mgo "gopkg.in/mgo.v2"
	"qiniu.com/app/common"
	log "qiniupkg.com/x/log.v7"
)

var db *mgo.Database

type daos struct {
	Job    JobDao
	Walker WalkerDao
}

var Daos *daos

func NewMongoSession(dbConf *common.MongoConfig) (session *mgo.Session, e error) {
	session, e = mgo.Dial(dbConf.Host)
	if e != nil {
		return
	}

	session.SetPoolLimit(dbConf.SessionPoolSize)

	db = session.DB(dbConf.DBName)
	return
}

func Init() (e error) {
	Daos = &daos{}
	j, e := NewJobDao(db)
	Daos.Job = j
	if e != nil {
		log.Warn("init job collection failed")
		return
	}

	w, e := NewWalkerDao(db)
	Daos.Walker = w
	if e != nil {
		log.Warn("init walker collection failed")
		return
	}

	return
}
