package store 

import (
	"log"
	mgo "gopkg.in/mgo.v2"
)
type MessageRepo struct {}
func(r MessageRepo)insertMessage(message Message)bool{
	session, err := mgo.Dial(SERVER1)
	defer session.Close()
	session.DB(DBNAME1).C(COLLECTION1).Insert(message)
	if err != nil {
		log.Fatalln("Insert", err)
		return false
	}
	return true

}