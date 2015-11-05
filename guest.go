package main

import "database/sql"

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
