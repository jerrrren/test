package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

/*
online postgres server
const (
	host     = "ec2-54-211-255-161.compute-1.amazonaws.com"
	port     =  5432
	user     = "fdikfgagiyyaxq"
	password = "e6398263fc3d8545d90cb455bc46cbcf7c47fd74ef9e6a8633e15c48c90feeff"
	dbname   = "db2jrt4m1j59dh"
)
*/

var DB *sql.DB = setupDatabase()

func setupDatabase() *sql.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		host     = os.Getenv("HOST")
		port, _  = strconv.Atoi(os.Getenv("DB_PORT"))
		user     = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
		dbname   = os.Getenv("DB")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", //need to change when uploading
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
