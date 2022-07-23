package forgetpassword

import (
	"github.com/gin-gonic/gin"
)

func ForgetPasswordRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("email/sentResetEmail", sentResetEmail())
	incomingRoutes.POST("resetPassword/reset",resetPassword())
}