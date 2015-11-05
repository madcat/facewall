package main

import "log"

func fatalWhenError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
