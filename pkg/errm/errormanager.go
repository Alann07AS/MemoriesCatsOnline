package errm

import (
	"errors"
	"fmt"
	"net/http"
)

func IfErrNotNil(err error, w http.ResponseWriter) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return true
	}
	return false
}

func ErrToJs(err error, w http.ResponseWriter) {
	http.Error(w, err.Error(), http.StatusBadRequest)
	fmt.Println(err)
}

var (
	ErrPlayerAlradyJoin = errors.New("go error: user already connected")
	ErrPlayerNoExist    = errors.New("go error: user no exist")
	ErrServerNoExist    = errors.New("go error: server no exist")
	ErrServerExist      = errors.New("go error: server already exist")
)
