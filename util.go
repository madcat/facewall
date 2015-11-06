package main

import (
	"log"
	"math/rand"
	"time"
)

func fatalWhenError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Shuffle(a []string) {
	rand.Seed(time.Now().UnixNano())
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
