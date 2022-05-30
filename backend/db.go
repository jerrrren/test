package main

import (
	"database/sql"
	"fmt"
)

const (
	host     = "localhost"
	port     =  5400
	user     = "postgres"
	password = "docker"
	dbname   = "postgres"
)

func setupDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")
	return db

}
