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
	AuthRoutes(router)
	UserRoutes(router)
	router.Run("intronus.herokuapp.com")
	defer db.Close()
}
