package dao

import (
	"github.com/globalsign/mgo"
)

//SeedDAO struct
type SeedDAO struct {
	*mgo.Session
}

//NewSeedDAO creating new seed dao object
func NewSeedDAO(conn string) (*SeedDAO, error) {
	s, err := mgo.Dial(conn)
	return &SeedDAO{
		Session: s,
	}, err
}

//GetSeed getting single seed object from mongo
func (handler *SeedDAO) GetColl() *mgo.Collection {
	s := handler.getFreshSession()
	//defer s.Close()
	return s.DB("seed").C("seeds")
}

func (handler *SeedDAO) getFreshSession() *mgo.Session {
	return handler.Session.Copy()
}
