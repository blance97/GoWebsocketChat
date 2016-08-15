package main

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	//"io"
	"log"
	"net/http"
	"os"
	"strings"
	//"time"
)

// Chat server.
type Server struct {
	Roomname  string
	messages  []Message
	clients   map[string]Client
	addCh     chan Client
	delCh     chan Client
	sendAllCh chan Message
	doneCh    chan bool
	errCh     chan error
}

// Create new chat server.
func NewServer(Roomname string) Server {
	messages := []Message{}
	clients := make(map[string]Client)
	addCh := make(chan Client) //inialize channlls
	delCh := make(chan Client)
	sendAllCh := make(chan Message)
	doneCh := make(chan bool)
	errCh := make(chan error)
	file := strings.Split(Roomname, "/")
	_, err := os.OpenFile("log/"+file[2], os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	// if _, err := io.WriteString(f, "Log File: "+Roomname+"\n "); err != nil {
	// 	log.Println("Log file Error:", err)
	// }
	return Server{
		Roomname,
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s Server) updateLog(text string) {
	file := strings.Split(s.Roomname, "/")
	f, err := os.OpenFile("log/"+file[2], os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(text + "\n"); err != nil {
		panic(err)
	}

}

/**
Store Client into add addclient channel
*/
func (s Server) Add(c Client) {
	s.addCh <- c
}

/**
Store Client into add delclient channel
*/
func (s Server) Del(c Client) {
	s.delCh <- c
}

/**
Store message into add message channel(type message)
*/
func (s Server) SendAll(msg Message) {
	s.sendAllCh <- msg
}

/**
Store boolean{true} into boolean channel
*/
func (s Server) Done() {
	s.doneCh <- true
}

/**
Store Error into err chanel
*/
func (s Server) Err(err error) {
	s.errCh <- err
}

/**
The server stores messages into s.messages.  writes out entire message history to client
*/
func (s Server) sendPastMessages(c Client) {
	file := strings.Split(s.Roomname, "/")
	stuff, err := readLines("log/" + file[2])
	if err != nil {
		log.Println("read: ", err)
	}
	for _, v := range stuff {
		var msg Message
		json.Unmarshal([]byte(v), &msg)
		c.Write(msg)
	}
}

/**
Sends message to all conected clients
*/
func (s Server) sendAll(msg Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s Server) Listen(RoomName string) {

	log.Println("Listening server...", RoomName)

	// websocket handler
	/**
	When the web socket is connected(created), it creates a new client. And starts the thread: client.Listen()
	*/
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()
		username, _ := getUsername(sid)
		client := NewClient(ws, s, username) //create new websocket with to server
		s.Add(client)                        //Add client to server
		client.Listen()                      //Fires go routine to listen write
	}
	http.Handle(s.Roomname, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh: // send to channel c
			log.Println("Added new client ", RoomName)
			s.clients[c.Username] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh: // Send to channgel C
			log.Println("Delete client")
			delete(s.clients, c.Username)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.updateLog(msg.String())
			//log.Println(reflect.TypeOf(msg.String()))
			//	StoreRoomInfo(s.Roomname, strings.Join(s.messages,","))
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
