package main

import (
	//"log"
	"encoding/json"
	"io/ioutil"
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

func SetSessionID(w http.ResponseWriter, r *http.Request) {
	// mu.Lock()
	// defer mu.Unlock()
	username := r.FormValue("Username")
	password := r.FormValue("password")
	log.Printf("Username %s pass %s", username, password)
	socketClientIP := strings.Split(r.RemoteAddr, ":")
	if r.URL.Path == "/login" {
		p, err := getUserPassword(username)
		if err != nil {
			log.Println("Error in getpassword ", err)
			http.Error(w, "Invalid Username or Password", 400)
			return
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
			sid = token
			StoreUserInfo(socketClientIP[0], username, password, sid)
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Invalid Username or Passwords", 400)
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	log.Println("Logout Hanlder")
	cookies, err := r.Cookie("SessionToken")
	if err != nil {
		log.Println("Error obtaining SessionToken", err)
		return
	}
	SessionToken := cookies.Value
	username, _ := getUsername(SessionToken)
	log.Println("USERNAME", username)
	log.Println("Cleared Logout")
	cookie := &http.Cookie{
		Name:  "SessionToken",
		Value: "0",
	}
	storeNewSessionToken(cookie.Value, username)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	return
}

/**
checks the SessionID
*/
func checkSession(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("SessionToken")
	SessionToken := cookie.Value
	sid = SessionToken
	log.Println("Update SessionToken")
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

func RoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data := getJSON(r)
		mu.Lock()
		server := NewServer("/entry/" + data["RoomName"].(string)) // start server
		mu.Unlock()
		go server.Listen(data["RoomName"].(string))
	}
}

type Roomier struct {
	Rooms []string
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	files, _ := ioutil.ReadDir("log/")
	var Room []string
	for _, file := range files {
		Room = append(Room, file.Name())
		//log.Println(file.Name())
	}
	q := Roomier{Rooms: Room}
	json.NewEncoder(w).Encode(q)
	//log.Println(len(files))
}

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
