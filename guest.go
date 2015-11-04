package main

type Guest struct {
	id     string
	code   string
	name   string
	group  string
	imgUrl string
}

type Win struct {
	order int32
	gid   string
	prize string
}
