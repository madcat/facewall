package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Guest struct {
	Gid    int
	Code   string
	Name   string
	Tag    string
	ImgUrl string
}

type Win struct {
	Wid   int
	Step  int
	Gid   string
	Prize string
}

type Winner struct {
	Step   int
	Code   string
	Name   string
	Tag    string
	ImgUrl string
	Prize  string
}

type Assignment struct {
	Tag    string
	Prize  string
	MaxWin int
}

type Controller struct {
	db *sql.DB
}

func (ctrl *Controller) MaxWinFor(tag string, prize string) int {
	sql := fmt.Sprintf("SELECT maxWin FROM assignment WHERE tag='%s' AND prize='%s'", tag, prize)
	var max int
	err := ctrl.db.QueryRow(sql).Scan(&max)
	if err != nil {
		return 0
	}
	return max
}

func (ctrl *Controller) InsertGuest(g Guest) (int64, error) {

	stmt, err := ctrl.db.Prepare("INSERT guest SET code=?,name=?,tag=?,imgUrl=?")
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(g.Code, g.Name, g.Tag, g.ImgUrl)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (ctrl *Controller) GetAllGuests(w http.ResponseWriter, r *http.Request) {
	rows, err := ctrl.db.Query("SELECT gid,code,name,tag,imgUrl FROM guest")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var guests []Guest
	for rows.Next() {
		var g Guest
		err = rows.Scan(&g.Gid, &g.Code, &g.Name, &g.Tag, &g.ImgUrl)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		guests = append(guests, g)
	}

	b, err := json.Marshal(guests)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(b)
}

func (ctrl *Controller) GetAllWinnners(w http.ResponseWriter, r *http.Request) {
	rows, err := ctrl.db.Query("SELECT win.step,code,name,tag,imgUrl,win.prize FROM guest JOIN win on win.gid=guest.gid")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var winners []Winner
	for rows.Next() {
		var winner Winner
		err = rows.Scan(&winner.Step, &winner.Code, &winner.Name, &winner.Tag, &winner.ImgUrl, &winner.Prize)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		winners = append(winners, winner)
	}

	b, err := json.Marshal(winners)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(b)
}
