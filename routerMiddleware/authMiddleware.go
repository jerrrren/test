package routerMiddleware

import(
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"

	"github.com/bojie/orbital/backend/auth"
	
)

func Authenticate() gin.HandlerFunc{
	return func(c*gin.Context){
		clientToken:=c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusInternalServerError,gin.H{"error":fmt.Sprintf("No Authorization Header provided")})
			c.Abort()
			return
		}	
		claims,err := auth.ValidateToken(clientToken)
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

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "https://intronusfrontend.herokuapp.com")
		//c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT,PATCH")

        if c.Request.Method == "OPTIONS" {	
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}