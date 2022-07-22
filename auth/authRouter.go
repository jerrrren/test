package auth

import (
	"github.com/gin-gonic/gin"

	
)




func AuthRoutes(incomingRoutes *gin.Engine){

	incomingRoutes.POST("users/signup",Signup())
	incomingRoutes.POST("users/login",Login())
	
	
}