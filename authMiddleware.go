package main

import(
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc{
	return func(c*gin.Context){
		clientToken:=c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusInternalServerError,gin.H{"error":fmt.Sprintf("No Authorization Header provided")})
			c.Abort()
			return
		}	
		claims,err := ValidateToken(clientToken)
		if err!=""{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			c.Abort()
			return
		}
		c.Set("name",claims.Name)
		c.Set("user_type",claims.User_type)
		c.Next()

	}
}