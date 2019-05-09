package store

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

type MessageRepo struct{}

func (r MessageRepo) insertMessage(message Message) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	session.DB(DBNAME).C(COLLECTION).Insert(message)
	if err != nil {
		log.Fatalln("Insert", err)
		return false
	}
	return true

}
