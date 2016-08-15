package main

import (
	//"log"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// Chat client.
type User struct {
	IP          string
	Username    string
	Password    string
	SessionID   string
	DateCreated int64
}

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
var sid string
func RoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		server := NewServer("/entry/room2") // start server
		go server.Listen("room2")
	}
}



func SetSessionID(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	socketClientIP := strings.Split(r.RemoteAddr, ":")
	cookie, err := r.Cookie("SessionToken")
	username := r.FormValue("Username")
	password := r.FormValue("password")
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "SessionToken",
			Value: "0",
		}
	}
	if r.URL.Path == "/login" {
		log.Println("login")
		redirectTarget := "/"
		p, err := getUserPassword(username)
		if err != nil {
			log.Println("Error in getpassword ", err)
			http.Redirect(w, r, "/login.html", 302)
		}

		var token string
		if password == p {
			log.Printf("Password1: %s Password2: %s", password, p)
			expiration := time.Now().Add(365 * 24 * time.Hour)
			for {
				token, _ = GenerateRandomString(64)
				if CheckValidSessionToken(token) {
					break
				}
			}
			cookie := &http.Cookie{Name: "SessionToken", Value: token, Expires: expiration}
			http.SetCookie(w, cookie)
			redirectTarget = "/chat.html"
			sid = token
			StoreUserInfo(socketClientIP[0], username, password, token)
		}
		http.Redirect(w, r, redirectTarget, 302)
		http.SetCookie(w, cookie)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	log.Println("Logout Hanlder")
	cookies, _ := r.Cookie("SessionToken")
	SessionToken := cookies.Value
	username, _ := getUsername(SessionToken)
	log.Println("Cleared Logout")
	cookie := &http.Cookie{
		Name:  "SessionToken",
		Value: "0",
	}
	http.SetCookie(w, cookie)
	storeNewSessionToken(cookie.Value, username)
}

/**
checks the SessionID
*/
func checkSession(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("SessionToken")
	SessionToken := cookie.Value
	if SessionToken != "0" {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "No Session", 403)
	return
}

/**
JSON Decoder
*/
func getJSON(r *http.Request) map[string]interface{} {
	var data map[string]interface{}

	//	log.Printf("getJSON:\tBegin execution")
	if r.Body == nil {
		log.Printf("No Request Body")
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("Error Decoding JSON")
	}
	defer r.Body.Close()
	return data
} //decode JSON

/**
Returns User so that it can validate whether or not a message belongs to them.
*/
func getUser(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("SessionToken")
	SessionToken := cookie.Value
	username, _ := getUsername(SessionToken)
	q := User{
		Username: username,
	}
	json.NewEncoder(w).Encode(q)
}

/**
Function to attain the users info and then store it in to the database
*/
func signUp(w http.ResponseWriter, r *http.Request) {
	//	log.Printf("Get User Handler")
	//defer log.Printf("done Get User Handler")
	socketClientIP := strings.Split(r.RemoteAddr, ":")
	data := getJSON(r)
	StoreUserInfo(socketClientIP[0], data["Username"].(string), data["Pass"].(string), "0")
}
