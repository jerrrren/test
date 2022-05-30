package main

import (
	"database/sql"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	db = setupDatabase()
	router := gin.Default()
	AuthRoutes(router)
	UserRoutes(router)
	router.Run()
	defer db.Close()
}
