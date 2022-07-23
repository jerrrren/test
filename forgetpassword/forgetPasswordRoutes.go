package forgetpassword

import (
	"github.com/gin-gonic/gin"
)

func ForgetPasswordRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("email/sentResetEmail", sentResetEmail())
	incomingRoutes.POST("resetPassword/reset/:token",resetPassword())
}