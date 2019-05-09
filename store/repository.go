package store

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

type Repository struct{}

const SERVER = "mongodb://susmita123:susmita123@ds145704.mlab.com:45704/newz-app-database"

// DBNAME the name of the DB instance
const DBNAME = "newz-app-database"

// COLLECTION is the name of the collection in DB
const COLLECTION = "chat-cache"

var productId = 10

// AddProduct adds a Product in the DB
func (r Repository) AddUser(product User) bool {
	log.Println(SERVER)
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	session.DB(DBNAME).C("users").Insert(product)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// AddComment adds a Comment in the DB
func (r Repository) AddCommentWithIndex(comment AddChat) bool {
	log.Println("add comment with  index")
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	session.DB("newz-app-database").C("messages").Insert(comment)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}
func (r Repository) AddCommentWithoutIndex(comment AddChat) bool {

	session, err := mgo.Dial(SERVER)
	defer session.Close()

	session.DB("testChat1").C("chat2").Insert(comment)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}
func (r Repository) AddNewsSection(comment AddChat) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()

	session.DB("testChat").C("parent-child-mapper").Insert(comment)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func (r Repository) GetChatHistoryFromDBWithIndex(comment AddChat) []AddChat {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	var ch []AddChat
	session.DB("testChat1").C("chat1").Find(bson.M{"parent": comment.Parent}).All(&ch)
	var m bson.M
	err = session.DB("testChat1").C("chat1").Find(bson.M{"parent": comment.Parent}).Explain(&m)
	if err == nil {
		fmt.Printf("Explain: %#v\n", m)
	}
	if err != nil {
		log.Fatal(err)
	}
	return ch
}
func (r Repository) GetChatHistoryFromDBWithoutIndex(comment AddChat) []AddChat {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	var ch []AddChat
	session.DB("testChat1").C("chat2").Find(bson.M{"parent": comment.Parent}).All(&ch)
	var m bson.M
	err = session.DB("testChat1").C("chat2").Find(bson.M{"parent": comment.Parent}).Explain(&m)
	if err == nil {
		fmt.Printf("Explain: %#v\n", m)
	}
	if err != nil {
		log.Fatal(err)
	}
	return ch
}
