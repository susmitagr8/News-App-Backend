package store

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	mgo "gopkg.in/mgo.v2"
)

//Controller ...
type Controller struct {
	Repository Repository
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
	_ = json.NewDecoder(req.Body).Decode(&user)
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
