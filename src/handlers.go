package main

import (
	//"log"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
type Rooms struct {
	Owner       string
	Roomname    string
	Private     string
	Password    string
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
	cookie, err := r.Cookie("SessionToken")
	if err != nil {
		log.Println("Error obtaining sid: ", err)
	}
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

type Roomier struct {
	Rooms []string
	Private []bool
}

func RoomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		data := getJSON(r)
		log.Println(data)
		_, err := os.OpenFile("log/"+data["Roomname"].(string), os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
		}
		StoreRoomInfo(data["Owner"].(string), data["Roomname"].(string), data["Private"].(string), data["RoomPass"].(string))
		server := NewServer("/entry/" + data["Roomname"].(string)) // start server
		go server.Listen(data["Roomname"].(string))
		w.WriteHeader(http.StatusOK)

	case "GET":
		mu.Lock()
		defer mu.Unlock()
		files, _ := ioutil.ReadDir("log/")
		var Room []string
		var Priv []bool
		for _, file := range files {
			Room = append(Room, file.Name())
			Priv = append(Priv, PrivateRoomChecker(file.Name()))
			//log.Println(file.Name())
		}
		q := Roomier{Rooms: Room,Private:Priv}
		json.NewEncoder(w).Encode(q)
		//log.Println(len(files))
	}
}
func DeleteRoom(w http.ResponseWriter, r *http.Request){

	data := r.URL.Query()
	Roomname := data.Get("RoomName")
	server := NewServer("/entry/" + Roomname)
		log.Println(server)
Done()//test for now
	w.WriteHeader(http.StatusOK)
	return
}
func CheckPrivateRoom(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	Roomname := data.Get("RoomName")
	log.Println("CheckPrivateRoom Called")
	if PrivateRoomChecker(Roomname) {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Not Private Room", 403)
	return
}
func CheckRoomPass(w http.ResponseWriter, r *http.Request){
	data:=getJSON(r)
	if data["RoomPass"].(string) == GetPrivateRoomPass(data["RoomName"].(string)){
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Wrong Password", 403)
	return
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
