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
	router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://intronus-frontend.herokuapp.com/"},
    AllowMethods:     []string{"PUT", "PATCH","GET","POST"},
    AllowHeaders:     []string{"Origin"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    AllowOriginFunc: func(origin string) bool {
      return origin == "https://github.com"
    },
    MaxAge: 12 * time.Hour,
    }))
	AuthRoutes(router)
	UserRoutes(router)
	router.Run()
	defer db.Close()
}
