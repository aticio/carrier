package dao

import (
	m "carrier/Seed/domestic/model"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
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

//AddSeed adding new seed to mongo
func (handler *SeedDAO) AddSeed(seed m.Seed) error {
	s := handler.getFreshSession()
	defer s.Close()
	err := s.DB("seed").C("seeds").Insert(&seed)
	if err != nil {
		log.Errorf("Error inserting soseedurce to mongo")
		return err
	}
	return nil
}

//GetSeed getting single seed object from mongo
func (handler *SeedDAO) GetSeed(host string, port string) (m.Seed, error) {
	s := handler.getFreshSession()
	defer s.Close()
	seed := m.Seed{}
	err := s.DB("seed").C("seeds").Find(bson.M{"host": host, "port": port}).One(&seed)
	return seed, err
}

//GetAll getting all seed objects from mongo
func (handler *SeedDAO) GetAll() ([]m.Seed, error) {
	s := handler.getFreshSession()
	defer s.Close()
	seeds := []m.Seed{}
	err := s.DB("seed").C("seeds").Find(nil).All(&seeds)
	return seeds, err
}

//UpdateSeed updates single Seed object by host and port
func (handler *SeedDAO) UpdateSeed(seed m.Seed) error {
	s := handler.getFreshSession()
	defer s.Close()
	err := s.DB("seed").C("seeds").Update(bson.M{"host": seed.Host, "port": seed.Port}, &seed)
	if err != nil {
		log.Errorf("Error updating seed %s,%s", seed.Host, seed.Port)
		return err
	}
	return nil
}

func (handler *SeedDAO) getFreshSession() *mgo.Session {
	return handler.Session.Copy()
}
