package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // sqlite

	"github.com/tunedmystic/commits.lol/app/config"
	"github.com/tunedmystic/commits.lol/app/db"
	"github.com/tunedmystic/commits.lol/app/server"
)

func main() {
	config := config.GetConfig()
	db := db.NewSqliteDB(config.DatabaseName)
	s := server.NewServer(db)

	addr := "localhost:8000"
	fmt.Printf("Running server on %v ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, s.Routes()))
}
