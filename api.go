package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

func main() {
	initDb()
	defer db.Close()

	//http.HandleFunc("/api/index", api)
	//http.HandleFunc("/api/repo/", handlers.)
	log.Fatal(http.ListenAndServe("localhost:8100", nil))
}

func initDb() {
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport], config[dbuser], config[dbpass], config[dbname])

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	//conf[dbhost] = "localhost"
	//conf[dbport] = "5432"
	//conf[dbuser] = "stellar"
	//conf[dbpass] = "123"
	//conf[dbname] = ""

	host, ok := os.LookupEnv(dbhost)
	if !ok {
		//panic("DBHOST environment variable required but not set")
		host = "localhost"
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		//panic("DBPORT environment variable required but not set")
		port = "5432"
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		//panic("DBUSER environment variable required but not set")
		user = "stellar"
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		//panic("DBPASS environment variable required but not set")
		password = "123"
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		//panic("DBNAME environment variable required but not set")
		name = "test"
	}

	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name

	return conf
}
