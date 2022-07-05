package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/lib/pq"


	"github.com/bojie/orbital/backend/db"
	"github.com/gin-gonic/gin"
)


type User struct {
	ID            uint   `json:"uid"`
	Name          string `json:"username"`
	Password      string `json:"password"`
	User_type     string `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Refresh_token string `json:"refresh_token"`
	Token         string `json:"token"`
}

func HashPassword(password string) string{
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err!=nil{
		fmt.Println(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err!= nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check=false
	}
	
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//validationErr := validate.Struct(user) //check detail of user struct variable
		//if validationErr != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		//return
		//}
		//check duplicate rows not included
		token, refreshToken, _ := GenerateAllTokens(user.Name, user.User_type)
		user.Token = token
		user.Refresh_token = refreshToken

		password := HashPassword(user.Password);
		user.Password = password

		result, err := db.DB.Exec("INSERT INTO users (name,password,refresh_token,token,user_type) VALUES ($1, $2, $3,$4,$5)", user.Name, user.Password, user.Refresh_token, user.Token, user.User_type)

		if err != nil {
			if error_code, ok := err.(*pq.Error); ok {
				if(error_code.Code == "23505"){
					c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "This username is already in use, please choose another one"})
					return
				}
			}

			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, result)
	}
}



func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		var foundUser User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		row := db.DB.QueryRow("SELECT * FROM users WHERE (name = $1)", user.Name)

		if err := row.Scan(&foundUser.ID, &foundUser.Name, &foundUser.Password, &foundUser.Token, &foundUser.Refresh_token, &foundUser.User_type); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusBadRequest, gin.H{"message": "password or username incorrect"})
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

			check_password,_ := VerifyPassword(foundUser.Password,user.Password);
			if(!check_password){
				c.JSON(http.StatusBadRequest, gin.H{"message": "password or username incorrect"})
			}

			return
		}

		token, refreshToken, err := GenerateAllTokens(foundUser.Name, foundUser.User_type)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}


		UpdateAllTokens(token, refreshToken, foundUser.ID)


		newrow := db.DB.QueryRow("SELECT * FROM users WHERE (name = $1)", user.Name)

		if err := newrow.Scan(&foundUser.ID, &foundUser.Name, &foundUser.Password, &foundUser.Token, &foundUser.Refresh_token, &foundUser.User_type); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusBadRequest, gin.H{"message": "password or username incorrect error"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}



		type response struct{
				ID            uint   `json:"uid"`
				Name          string `json:"username"`
				User_type     string `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
				Refresh_token string `json:"refresh_token"`
				Token         string `json:"token"`
		}

		var json_response response;
		json_response.Name = foundUser.Name
		json_response.ID = foundUser.ID
		json_response.Refresh_token = foundUser.Refresh_token
		json_response.Token = foundUser.Token
		json_response.User_type = foundUser.User_type

		c.JSON(http.StatusOK, json_response)
	}

}




func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
