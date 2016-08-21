package main

import (
	"encoding/json"
	"log"
)

type Message struct {
	Roomname string `json:"Roomname"`
	Author   string `json:"author"`
	Time     string `json:"time"`
	Body     string `json:"body"`
}

func (self *Message) String() string {
	msg := &Message{
		Roomname: self.Roomname,
		Author:   self.Author,
		Time:     self.Time,
		Body:     self.Body,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(b)
}
