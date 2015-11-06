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
	Prize  string
	Step   int
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
	rows, err := ctrl.db.Query("SELECT step,prize,code,name,tag,imgUrl FROM guest WHERE prize IS NOT NULL AND prize != '' ORDER BY step")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var winners []Guest
	for rows.Next() {
		var winner Guest
		err = rows.Scan(&winner.Step, &winner.Prize, &winner.Code, &winner.Name, &winner.Tag, &winner.ImgUrl)
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
