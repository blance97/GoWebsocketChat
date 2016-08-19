package main

import (
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var db = InitDB("database/ChatDB")

func main() {
	CreateUserTable()
	createRoomTable()
	initRooms()
	log.SetFlags(log.Lshortfile)
	//	server := NewServer("/entry/room1")
	//go server.Listen("room1")
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	http.HandleFunc("/createRoom", RoomHandler)
	http.HandleFunc("/deleteRoom/", DeleteRoom)
	http.HandleFunc("/getRooms", RoomHandler)
	http.HandleFunc("/RoomPerm/",CheckPrivateRoom)
	http.HandleFunc("/RoomPassCheck",CheckRoomPass)
	http.HandleFunc("/login", SetSessionID)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/getUser", getUser)
	http.HandleFunc("/checkSession", checkSession)
	//	log.Printf("Running on port %d\n", *port)
	log.Fatal(http.ListenAndServe(":80", nil))
}
