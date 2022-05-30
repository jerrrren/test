package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	db = setupDatabase()
	router := gin.Default()
	router.Use(CORSMiddleware())
	AuthRoutes(router)
	UserRoutes(router)
	router.Run()
	defer db.Close()
}
