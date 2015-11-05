package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
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

func configFromFlags() *DBConfig {
	var dbconfig DBConfig
	flag.StringVar(&dbconfig.url, "dburl", "", "host:port for database")
	flag.StringVar(&dbconfig.name, "dbname", "", "database name")
	flag.StringVar(&dbconfig.user, "dbuser", "", "database user")
	flag.StringVar(&dbconfig.pass, "dbpass", "", "database password")
	flag.Parse()
	if dbconfig.url == "" || dbconfig.name == "" || dbconfig.user == "" || dbconfig.pass == "" {
		return nil
	}
	return &dbconfig
}

func fatalWhenError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbconfig := configFromFlags()
	if dbconfig == nil {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	db, err := sql.Open("mysql", dbconfig.connectionString())
	fatalWhenError(err)
	fatalWhenError(db.Ping())
}
