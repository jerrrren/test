package main

import (
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/users",Authenticate(),GetUsers())
	incomingRoutes.GET("/users/:user_id",Authenticate(),GetUser())
}			