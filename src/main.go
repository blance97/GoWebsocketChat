package main

import (
	"log"
	"net/http"
)

var db = InitDB("database/ChatDB")

func main() {
	CreateUserTable()
	log.SetFlags(log.Lshortfile)
	server := NewServer("/entry/room1")
	go server.Listen("room1")
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	 http.HandleFunc("/createRoom", RoomHandler)
	 http.HandleFunc("/login", SetSessionID)
	 http.HandleFunc("/logout", logout)
	 http.HandleFunc("/signup", signUp)
	 http.HandleFunc("/getUser", getUser)
	 http.HandleFunc("/checkSession", checkSession)
	//	log.Printf("Running on port %d\n", *port)
	log.Fatal(http.ListenAndServe(":80", nil))
}