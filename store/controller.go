package store

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Controller ...
type Controller struct {
	Repository  Repository
	MessageRepo MessageRepo
}

/* Middleware handler to handle all requests for authentication */
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					log.Println("TOKEN WAS VALID")
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

// Get Authentication token GET /
func (c *Controller) GetToken(w http.ResponseWriter, req *http.Request) {
	var user RegisteredUser
	v := json.NewDecoder(req.Body).Decode(&user)
	if v != nil {
		fmt.Println(v)
	}
	fmt.Println(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})

	log.Println("Username: " + user.Username)
	log.Println("Password: " + user.Password)
	result := User{}
	session, err := mgo.Dial("localhost:27017")
	defer session.Close()

	c1 := session.DB("news-user").C("store")
	err = c1.Find(bson.M{"username": user.Username}).One(&result)
	if err != nil {
		errState := ResponseMessageError{"Not correct Credentials"}
		w.WriteHeader(http.StatusUnauthorized)
		succ, _ := json.Marshal(errState)
		w.Write(succ)
		return
	}
	if user.Password != result.Password {
		errState := ResponseMessageError{"Not correct Credentials"}
		w.WriteHeader(http.StatusUnauthorized)
		succ, _ := json.Marshal(errState)
		w.Write(succ)
		return
	}
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

//Add user after registration /
func (c *Controller) AddUser(w http.ResponseWriter, r *http.Request) {
	log.Println("in add User handle")
	var product User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // read the body of the request
	log.Println(body)
	if err != nil {
		log.Fatalln("Error AddUser", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddUser", err)
	}
	if err := json.Unmarshal(body, &product); err != nil { // unmarshall body contents as a type Candidate
		log.Println(err)
	}

	log.Println(product)
	success := c.Repository.AddUser(product) // adds the user to the DB
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	responseMsg := ResponseMessage{"Successful addition of user"}
	succ, _ := json.Marshal(responseMsg)
	w.Write(succ)
}

//Add a message by POST
func (c *Controller) AddMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(body)
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error addMessage", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &message); err != nil {
		log.Println(err)
	}
	log.Println(message)
	success := c.MessageRepo.insertMessage(message)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	responseMsg := ResponseMessage{"Successful addition of message"}
	succ, _ := json.Marshal(responseMsg)
	w.Write(succ)
}

//Retrieve messages by GET

func (c *Controller) FindMessage(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	log.Println(r.URL.Query())
	pages, ok := r.URL.Query()["PageNumber"]
	if !ok || len(pages[0]) < 1 {
		session, err := mgo.Dial(SERVER1)
		if err != nil {
			log.Fatalln(err)
		}
		if err := session.DB(DBNAME1).C(COLLECTION1).Find(nil).All(&messages); err != nil {
			log.Fatalln("Find error ", err)
		}
		defer session.Close()
		for _, message := range messages {
			author := message.Author_id
			mess := message.Message
			json.NewEncoder(w).Encode("Author :- " + author)
			json.NewEncoder(w).Encode("Message :- " + mess)
		}
		return
	}
	artid, ok := r.URL.Query()["ArticleID"]
	if !ok || len(artid[0]) < 1 {
		log.Fatalln("INVALID QUERY")
	}
	session, err := mgo.Dial(SERVER1)
	if err != nil {
		log.Fatalln(err)
	}
	n := pages[0]
	num, err := strconv.Atoi(n)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	q := session.DB(DBNAME1).C(COLLECTION1).Find( bson.M{"article_id": artid[0]}).Limit(PageSize)
	q = q.Skip((num - 1) * PageSize)
	if err := q.All(&messages); err != nil {
		log.Fatalln("Query error", err)
	}

	defer session.Close()
	log.Println("Successful query")
	//res:=ResData{messages}
	log.Println(len(messages))
	for _, message := range messages {
		author := message.Author_id
		mess := message.Message
		json.NewEncoder(w).Encode("Author :- " + author)
		json.NewEncoder(w).Encode("Message :- " + mess)
	}
	return
}
