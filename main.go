package main

import (
	"github.com/bojie/orbital/backend/auth"
	"github.com/bojie/orbital/backend/chat"
	"github.com/bojie/orbital/backend/db"
	"github.com/bojie/orbital/backend/email"
	"github.com/bojie/orbital/backend/pairing"
	"github.com/bojie/orbital/backend/routerMiddleware"
	"github.com/bojie/orbital/backend/user"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	wsServer := chat.NewWebSocketServer()
	go wsServer.Run()

	router := gin.Default()

	router.Use(routerMiddleware.CORSMiddleware())
	auth.AuthRoutes(router)
	user.UserRoutes(router)
	email.EmailRoutes(router)
	pairing.PairingRoutes(router)

	router.GET("/ws", chat.ServeWs(wsServer))
	router.Run("")

	defer db.DB.Close()
}
