package main

import (
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxId int = 0

// Chat client.
type Client struct {
	//Username   string
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan Message
	doneCh chan bool
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}
	maxId++
	ch := make(chan Message, channelBufSize)
	doneCh := make(chan bool)
	return Client{maxId, ws, server, ch, doneCh}
}

func (c Client) Write(msg Message) {
	c.ch <- msg //  If client receives a message do nothing)  This puts message into c.ch which is used to listen write
}

// Listen Write and Read request via chanel
func (c Client) Listen() {
	go c.listenWrite()
	c.listenRead() // while loop that is constanly listening for data and sending it to all clients.
}

// Listen write request via chanel
func (c Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {
		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh: // tells the server to delete current client( if there is new data?)
			c.server.Del(c)
			return
		}
	}
}

// Listen read request via chanel
/**
This function takes in a message from the client and then shares it with all the clients
*/
func (c Client) listenRead() {
	log.Println("Listening read from client")
	for {
		//Sends message to all clients(this one is always called )
		var msg Message
		err := websocket.JSON.Receive(c.ws, &msg)
		log.Println("MESSAGE: ", &msg)
		if err == io.EOF {
			c.doneCh <- true // if end of file delete user which is called in listnwrite
		} else if err != nil {
			c.server.Err(err)
		} else {
			c.server.SendAll(msg)
		}
	}

}
