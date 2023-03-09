package server

import (
	"fmt"
	"time"

	"test2/pkg/game"
)

type Server struct {
	Name       string           `json:"name"`
	Users      map[string]*User `json:"users"`
	Broadcast  chan []byte      `json:"-"`
	Register   chan *User       `json:"-"`
	Unregister chan *User       `json:"-"`
	Game       *game.MemoryGame `json:"-"`
}

func (s *Server) GetUserByName(username string) *User {
	return s.Users[username]
}

func NewServer(servername string) *Server {
	return &Server{Name: servername, Users: map[string]*User{}, Broadcast: make(chan []byte), Register: make(chan *User), Unregister: make(chan *User), Game: game.NewGame(servername)}
}

func (s *Server) Run() {
	for {
		select {
		case ic := <-s.Broadcast:
			inComing := MessageParse(ic, s)
			fmt.Println("INCOMING INSTRUCTION", inComing)
			actionsGO[inComing.Instruction](inComing.User, inComing.Params)
		case user := <-s.Register:
			s.Users[user.Name] = user
			go func() {
				time.Sleep(500 * time.Millisecond)
				s.SendToEveryOne(NewMessageOut(JS_UPDATE_PLAYER_READY, user.Server.Game.Players))
			}()
		case user := <-s.Unregister:
			fmt.Println(user.Name, "C'est deco")
			s.Game.Players[user.Name].Ready = false
			s.SendToEveryOne(NewMessageOut(JS_UPDATE_PLAYER_READY, s.Game.Players))
			go func() {
				// deco au bout de 30 seconde d'interuption
				time.Sleep(30 * time.Second)
				if _, ok := s.Users[user.Name]; !ok {
					s.Game.RemovePlayer(user.Name)
					s.SendToEveryOne(NewMessageOut(JS_UPDATE_PLAYER_READY, s.Game.Players))
				}
			}()
			close(user.Send)
			delete(s.Users, user.Name)
		}
	}
}

func (s *Server) UpDateGame() {
	for _, u := range s.Users {
		u.Send <- s.Game
	}
}

func (s *Server) StartGame() {
	go func() {
		// start after 10s
		tnow := time.Now()
		for {
			if time.Since(tnow) >= time.Second*10 {
				break
			}
			time.Sleep(100 * time.Millisecond)
			select {
			case <-s.Game.CancelStartChanel: // cancel start
				return
			default:
				continue
			}
		}

		// randomize play order
		s.Game.SetPlayerOrderRandom()
		s.Game.ActivePlayerIndex = 0
		s.SendToEveryOne(NewMessageOut(JS_SHOW_GAME, s.Game))
		time.Sleep(300 * time.Millisecond)
		s.SendToEveryOne(NewMessageOut(JS_UPDATE_PLIST, s.Game.Players, s.Game.ActivePlayerIndex))
		s.SendToPlayerName(s.Game.PlayerPlayDirection[0], NewMessageOut(JS_YOUR_TURN))
	}()
}

func (s *Server) SendToEveryOne(m MessageOut) {
	fmt.Println("MESSAGE SEND ___ ", m)
	for _, u := range s.Users {
		u.Send <- m
	}
}

func (s *Server) SendToPlayerName(name string, m MessageOut) {
	fmt.Println("MESSAGE SEND TO"+name+"___ ", m)
	s.Users[name].Send <- m
}
