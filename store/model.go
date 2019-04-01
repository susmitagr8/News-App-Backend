package store

type RegisteredUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Message struct {
	Article_id string `json:"art_id"`
	Author_id  string `json:"aut_id"`
	Message    string `json:"message"`
	Time string `json:"time"`
}
type ResData struct{
	Data []Message 
}
type JwtToken struct {
	Token string `json:"token"`                        
}

type Exception struct {
	Message string `json:"message"`
}
type ResponseMessage struct {
	Status string
}
type ResponseMessageError struct {
	Status string
}

// User represents an user who is registering
type User struct {
	Username string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
