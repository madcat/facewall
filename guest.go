package main

import (
	"database/sql"
	"encoding/json"
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

type Controller struct {
	db *sql.DB
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
