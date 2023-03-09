package server

import (
	"encoding/json"
	"fmt"
	"time"
)

// message entrant executer
type MessageIn struct {
	User        *User
	Instruction Order
	Params      []interface{}
}

func MessageParse(input []byte, s *Server) *MessageIn {
	// resoit un json format UserName, Instruction, [params]
	inputData := &struct {
		UserName    string
		Instruction int
		Params      []interface{}
	}{}
	json.Unmarshal(input, inputData)
	fmt.Println(string(input))
	return &MessageIn{
		User:        s.GetUserByName(inputData.UserName),
		Instruction: Order(inputData.Instruction),
		Params:      inputData.Params,
	}
}

type Order int

type MessageOut struct {
	Instruction Order
	Params      interface{}
}

func NewMessageOut(i Order, params ...interface{}) MessageOut {
	return MessageOut{i, params}
}

// orde a envoyer au js pour qu'il execute
const (
	JS_UPDATE_PLAYER_READY = 1
	JS_TOGGLE_START_CHRONO = 2
	JS_SHOW_GAME           = 3
	JS_SHOW_CARD           = 4
	JS_YOUR_TURN           = 5
	JS_UPDATE_PLIST        = 6
	JS_HIDING_CARD         = 7
)

// definition des action possible

var actionsGO = map[Order]func(u *User, args []interface{}){}

// ordre recu par le js pour etre executer ici
const (
	GO_TOGGLE_READY    = Order(1)
	GO_SEND_CARD_BY_ID = Order(2)
	GO_CHECK_TWIN      = Order(3)
	GO_4               = Order(4)
	GO_5               = Order(5)
	GO_6               = Order(6)
	GO_7               = Order(7)
	GO_8               = Order(8)
	GO_9               = Order(9)
	GO_10              = Order(10)
	GO_11              = Order(11)
	GO_12              = Order(12)
	GO_13              = Order(13)
	GO_14              = Order(14)
	GO_15              = Order(15)
	GO_16              = Order(16)
	GO_17              = Order(17)
	GO_18              = Order(18)
	GO_19              = Order(19)
	GO_20              = Order(20)
	GO_21              = Order(21)
	GO_22              = Order(22)
	GO_23              = Order(23)
	GO_24              = Order(24)
	GO_25              = Order(25)
	GO_26              = Order(26)
	GO_27              = Order(27)
	GO_28              = Order(28)
	GO_29              = Order(29)
	GO_30              = Order(30)
)

func init() {
	actionsGO[GO_TOGGLE_READY] = func(u *User, i []interface{}) {
		game := u.Server.Game
		wasReady := game.IsAllPlayersReady()
		game.GetPlayerByName(u.Name).ToggleReady()
		u.Server.SendToEveryOne(NewMessageOut(JS_UPDATE_PLAYER_READY, game.Players))
		if b := game.IsAllPlayersReady(); b || wasReady {
			u.Server.SendToEveryOne(NewMessageOut(JS_TOGGLE_START_CHRONO))
			if !b {
				fmt.Println("TRIGERCANCEL")
				game.CancelStartChanel <- struct{}{}
				fmt.Println("TRIGERCANCEL__valid")
			} else {
				u.Server.StartGame()
			}
		}
	}

	actionsGO[GO_SEND_CARD_BY_ID] = func(u *User, i []interface{}) {
		game, id := u.Server.Game, int(i[0].(float64))
		u.Server.SendToEveryOne(NewMessageOut(JS_SHOW_CARD, id, game.Cards[id-1].Img))
	}

	actionsGO[GO_CHECK_TWIN] = func(u *User, i []interface{}) {
		c1, c2 := int(i[0].(float64)), int(i[1].(float64))
		game := u.Server.Game
		s := u.Server
		if game.CheckTwin(c1, c2) {
			// good / false answer

			game.Players[u.Name].Score += 1
		} else {
			time.Sleep(2 * time.Second)
			s.SendToEveryOne(NewMessageOut(JS_HIDING_CARD, c1, c2))
			game.NextPlayerTurn()
		}
		s.SendToEveryOne(NewMessageOut(JS_UPDATE_PLIST, s.Game.Players, s.Game.ActivePlayerIndex))
		s.SendToPlayerName(game.GetPlayerTurnName(), NewMessageOut(JS_YOUR_TURN))
	}
}
