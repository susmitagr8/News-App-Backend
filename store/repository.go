package store

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

type Repository struct{}


var productId = 10

// AddProduct adds a Product in the DB
func (r Repository) AddUser(product User) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	session.DB(DBNAME).C(COLLECTION).Insert(product)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
