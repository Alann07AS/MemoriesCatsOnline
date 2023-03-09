package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type User struct {
	Name   string           `json:"name"`
	Conn   *websocket.Conn  `json:"-"`
	Send   chan interface{} `json:"-"`
	Server *Server          `json:"-"`
}

func NewUser(name string, conn *websocket.Conn, s *Server) *User {
	return &User{
		Name:   name,
		Conn:   conn,
		Send:   make(chan interface{}),
		Server: s,
	}
}

func (u *User) Read() {
	defer func() {
		u.Conn.Close()
		u.Server.Unregister <- u
	}()
	for {
		_, message, err := u.Conn.ReadMessage()
		if err != nil {
			fmt.Println("u.Server.Users", err)
			break
		}
		u.Server.Broadcast <- message
	}
}

func (u *User) Write() {
	defer u.Conn.Close()
	for {
		m, ok := <-u.Send
		if !ok {
			return
		}
		err := u.Conn.WriteJSON(m)
		if err != nil {
			fmt.Println("l:52, USER WrieteJSON ERR", err)
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(s *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	name := r.URL.Query().Get("username")
	client := &User{Name: name, Server: s, Conn: conn, Send: make(chan interface{})}
	client.Server.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.Write()
	go client.Read()
}
