package pairing

import (
	"github.com/gin-gonic/gin"
)

func PairingRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/pairing/registerAdd", AddRegistrationUser())
	incomingRoutes.POST("/pairing/ifPaired", checkIfPaired())
	incomingRoutes.POST("/pairing/fillAndMatch", FillAndMatch())
}
