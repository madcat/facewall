package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

/*
  shuffle states:
  - all
  - shuffling/{prize}
  - history/{step}
*/
type ShuffleController struct {
	state  string
	db     *sql.DB
	router *mux.Router
}

func NewShuffleController(r *mux.Router, db *sql.DB) *ShuffleController {
	sc := new(ShuffleController)
	sc.router = r
	sc.router.StrictSlash(true)
	sc.state = "all"
	sc.db = db

	sc.router.HandleFunc("/start/{prize}", sc.start).Methods("POST")
	sc.router.HandleFunc("/end", sc.end).Methods("POST")
	sc.router.HandleFunc("/history/{step}", sc.history).Methods("GET")
	sc.router.HandleFunc("/all", sc.all).Methods("GET")
	sc.router.HandleFunc("/state", sc.getState).Methods("GET")
	return sc
}

func (sc *ShuffleController) start(w http.ResponseWriter, r *http.Request) {
	sc.state = "shuffling/一等奖"
	w.Write([]byte(sc.state))
}

func (sc *ShuffleController) end(w http.ResponseWriter, r *http.Request) {
	sc.state = "history/2"
	w.Write([]byte(sc.state))
}

func (sc *ShuffleController) history(w http.ResponseWriter, r *http.Request) {
	sc.state = "history/1"
	w.Write([]byte(sc.state))
}

func (sc *ShuffleController) all(w http.ResponseWriter, r *http.Request) {
	sc.state = "all"
	w.Write([]byte(sc.state))
}

func (sc *ShuffleController) getState(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(sc.state))
}
