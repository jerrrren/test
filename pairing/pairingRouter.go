package pairing

import (
	"github.com/gin-gonic/gin"
)

func PairingRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/pairing/registerAdd", AddRegistrationUser())
	incomingRoutes.POST("/pairing/checkInfo", checkInfo())
	incomingRoutes.POST("/pairing/fillIndicators", FillIndicators())
	incomingRoutes.POST("/pairing/match", Match())
	incomingRoutes.POST("/pairing/addPairedUser", AddPaireduser())
	incomingRoutes.POST("/pairing/deleteSingleUser", DeleteSingleUser())
	incomingRoutes.POST("/pairing/deletePairedUser", DeletePairedUser())
}
