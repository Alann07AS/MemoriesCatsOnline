package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	"test2/pkg/errm"
	"test2/pkg/server"

	"github.com/gorilla/websocket"
)

var (
	serverList = make(map[string]*server.Server)
	serverPop  = make(map[string]int)
	users      = []string{}
	lastupdate = time.Now()
	upgrader   = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func init() {
	rand.Seed(time.Now().UnixNano()) // randomise la seed de base
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // HOME PAGE
		renderTmpl(w, "accueil", nil)
	})

	http.HandleFunc("/checkforupdate", func(w http.ResponseWriter, r *http.Request) { // SEND LAST UPDATE TIME
		for i := range serverList {
			if len(serverList[i].Users) != serverPop[i] {
				lastupdate = time.Now() // update last update
				serverPop[i] = len(serverList[i].Users)
			}
		}
		fmt.Fprint(w, lastupdate.String())
	})

	http.HandleFunc("/getservelist", func(w http.ResponseWriter, r *http.Request) { // SEND JSON SERVERLIST
		json, err := json.Marshal(serverList)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(json))
	})

	http.HandleFunc("/newserver", func(w http.ResponseWriter, r *http.Request) { // CREATE NEW SERVER
		name := r.URL.Query().Get("servername")
		if _, ok := serverList[name]; ok { // dont create server if alrady exist
			errm.ErrToJs(errm.ErrServerExist, w)
			return
		}
		s := server.NewServer(name) // create new server
		serverList[name] = s        // save new server
		serverPop[name] = 0         // save new server o$population
		go s.Run()
		lastupdate = time.Now() // update last update
	})

	http.HandleFunc("/deleteserver", func(w http.ResponseWriter, r *http.Request) { // DELETE SERVER
		name := r.URL.Query().Get("servername")
		if _, ok := serverList[name]; !ok { // dont create server if alrady exist
			errm.ErrToJs(errm.ErrServerNoExist, w)
			return
		}
		delete(serverList, name)
		delete(serverPop, name)
		lastupdate = time.Now() // update last update
	})

	http.HandleFunc("/loadgame", func(w http.ResponseWriter, r *http.Request) { // LOAD GAME PAGE
		renderTmpl(w, "game", nil)
	})

	http.HandleFunc("/joingame", func(w http.ResponseWriter, r *http.Request) { // JOIN a SERVER
		conn, err := upgrader.Upgrade(w, r, nil)
		if errm.IfErrNotNil(err, w) {
			return
		}
		sname := r.URL.Query().Get("servername")
		s, ok := serverList[sname]
		if !ok {
			errm.ErrToJs(errm.ErrServerNoExist, w)
			conn.Close()
			return
		}
		username, err := r.Cookie("username")
		if errm.IfErrNotNil(err, w) {
			conn.Close()
			return
		}
		_, ok = s.Users[username.Value]
		if ok {
			errm.ErrToJs(errm.ErrPlayerAlradyJoin, w)
			conn.Close()
			return
		}
		lastupdate = time.Now() // update last update

		client := server.NewUser(username.Value, conn, s)
		s.Register <- client
		client.Server.Game.AddPlayer(client.Name)

		go client.Write()
		go client.Read()
		serverPop[s.Name]++ // add pop
	})

	http.HandleFunc("/leavegame", func(w http.ResponseWriter, r *http.Request) { // LEAVE SERVER
		sname := r.URL.Query().Get("servername")
		s, ok := serverList[sname]
		if !ok {
			errm.ErrToJs(errm.ErrServerNoExist, w)
			return
		}
		uname, err := r.Cookie("username")
		if errm.IfErrNotNil(err, w) {
			return
		}
		u, ok := s.Users[uname.Value]
		if !ok {
			errm.ErrToJs(errm.ErrPlayerNoExist, w)
			return
		}
		u.Conn.Close()
		delete(s.Users, uname.Value)
		lastupdate = time.Now() // update last update
	})

	contain := func(s []string, str string) bool {
		if str == "" {
			return true
		}
		for _, v := range s {
			if v == str {
				return true
			}
		}
		return false
	}

	http.HandleFunc("/newuser", func(w http.ResponseWriter, r *http.Request) { // register newUser
		username := r.URL.Query().Get("username")

		fmt.Println("p", users, username)
		if contain(users, username) {
			fmt.Fprint(w, "Username already use")
			return
		}
		users = append(users, username)
	})

	http.HandleFunc("/checkuser", func(w http.ResponseWriter, r *http.Request) { // register newUser
		username := r.URL.Query().Get("username")
		if !contain(users, username) {
			fmt.Fprint(w, "Username not register")
		}
	})

	http.ListenAndServe(":8080", nil) // START MAIN SERVER
}

func renderTmpl(w http.ResponseWriter, tmplname string, data any) {
	tmpl, err := template.ParseFiles("template/" + tmplname + ".page.tmpl")
	if errm.IfErrNotNil(err, w) {
		return
	}
	tmpl.Execute(w, nil)
}
