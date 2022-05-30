package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "ec2-54-211-255-161.compute-1.amazonaws.com"
	port     =  5432
	user     = "fdikfgagiyyaxq"
	password = "e6398263fc3d8545d90cb455bc46cbcf7c47fd74ef9e6a8633e15c48c90feeff"
	dbname   = "db2jrt4m1j59dh"
)

func setupDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require",
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
