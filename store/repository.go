package store

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

//Repository ...
type Repository struct{}

// SERVER the DB server
const SERVER = "localhost:27017"

// DBNAME the name of the DB instance
const DBNAME = "news-user"

// COLLECTION is the name of the collection in DB
const COLLECTION = "store"

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
