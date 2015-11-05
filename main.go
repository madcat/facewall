package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type DBConfig struct {
	url  string
	name string
	user string
	pass string
}

func (config *DBConfig) connectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", config.user, config.pass, config.url, config.name)
}

func main() {
	var dbconfig DBConfig
	var port string
	flag.StringVar(&dbconfig.url, "dburl", "", "host:port for database")
	flag.StringVar(&dbconfig.name, "dbname", "", "database name")
	flag.StringVar(&dbconfig.user, "dbuser", "", "database user")
	flag.StringVar(&dbconfig.pass, "dbpass", "", "database password")
	flag.StringVar(&port, "port", "3000", "port server will run on")
	flag.Parse()
	if dbconfig.url == "" || dbconfig.name == "" || dbconfig.user == "" || dbconfig.pass == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	db, err := sql.Open("mysql", dbconfig.connectionString())
	fatalWhenError(err)
	fatalWhenError(db.Ping())

	var ctrl = Controller{db}

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/guest", ctrl.GetAllGuests).Methods("GET")
	router.HandleFunc("/win", ctrl.GetAllWinnners).Methods("GET")

	fmt.Printf("server running at %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
