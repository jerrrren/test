package email

import (
	"github.com/gin-gonic/gin"
)

func EmailRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("email/sendverificationemail", sendVerificationMessage())
	incomingRoutes.POST("email/verifyemail",verifyEmail())
	incomingRoutes.GET("email/checkverified/:user_id",getVerified())
}
