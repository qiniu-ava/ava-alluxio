package database

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"qiniu.com/app/common/typo"
)

type WalkerDao interface {
	InsertWalker(walker *typo.Walker) error
	UpdateWalker(walker *typo.Walker) error
	DeleteWalker(name string) error
	GetWalker(name string) (*typo.Walker, error)
	GetWalkers() ([]typo.Walker, error)
}

const walkerCollectionName = "walker"

type walkerDao struct {
	collection *mgo.Collection
}

func NewWalkerDao(db *mgo.Database) (WalkerDao, error) {
	c := db.C(walkerCollectionName)
	if e := c.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true, Name: "walker_name"}); e != nil {
		return nil, e
	}
	w := &walkerDao{
		collection: db.C(walkerCollectionName),
	}
	return w, nil
}

func (w *walkerDao) InsertWalker(walker *typo.Walker) error {
	// TODO fix createTime is 0 bug
	walker.CreateTime = time.Now()
	return w.collection.Insert(walker)
}

func (w *walkerDao) UpdateWalker(walker *typo.Walker) error {
	q := bson.M{"name": walker.Name}
	r := bson.M{}
	if walker.Jobs != nil {
		r["jobs"] = walker.Jobs
	}
	if &walker.Status != nil {
		r["status"] = walker.Status
	}
	if &walker.CreateTime != nil {
		r["createTime"] = walker.CreateTime
	}
	return w.collection.Update(q, bson.M{"$set": r})
}

func (w *walkerDao) GetWalker(name string) (wr *typo.Walker, e error) {
	q := bson.M{"name": name}
	r := w.collection.Find(q)
	c, e := r.Count()
	if c == 0 {
		return nil, errors.Errorf("not walker named %s found", name)
	} else if e != nil {
		return nil, e
	}

	wr = &typo.Walker{}
	r.One(wr)

	return wr, nil
}

func (w *walkerDao) DeleteWalker(name string) error {
	return w.collection.Remove(bson.M{"name": name})
}

func (w *walkerDao) GetWalkers() ([]typo.Walker, error) {
	walkers := []typo.Walker{}
	r := w.collection.Find(bson.M{"status": int(typo.OnLine)})
	iter := r.Iter()
	item := typo.Walker{}
	for iter.Next(&item) {
		walkers = append(walkers, item)
	}

	return walkers, nil
}
