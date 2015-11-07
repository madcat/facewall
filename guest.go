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
	rows, err := ctrl.db.Query("SELECT gid,code,name,tag,imgUrl,prize,step FROM guest")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var guests []Guest
	for rows.Next() {
		var g Guest
		var prize sql.NullString
		var step sql.NullInt64
		err = rows.Scan(&g.Gid, &g.Code, &g.Name, &g.Tag, &g.ImgUrl, &prize, &step)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if step.Valid {
			g.Step = int(step.Int64)
		}
		if prize.Valid {
			g.Prize = prize.String
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
	winners := make([]Guest, 0)
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

func (ctrl *Controller) GetAllPrizes(w http.ResponseWriter, r *http.Request) {
	rows, err := ctrl.db.Query("SELECT id,prize FROM prize ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var prizes []Prize
	for rows.Next() {
		var p Prize
		err = rows.Scan(&p.Id, &p.Prize)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		prizes = append(prizes, p)
	}

	b, err := json.Marshal(prizes)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(b)
}
