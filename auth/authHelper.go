package auth

import (
	"errors"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/bojie/orbital/backend/db"

)


type User struct {
	ID            uint   `json:"uid"`
	Name          string `json:"username"`
	Password      string `json:"password"`
	User_type     string `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Refresh_token string `json:"refresh_token"`
	Token         string `json:"token"`
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

		result, err := db.DB.Exec("INSERT INTO users (name,password,refresh_token,token,user_type) VALUES ($1, $2, $3,$4,$5)", user.Name, user.Password, user.Refresh_token, user.Token, user.User_type)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"test": err})
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
		row := db.DB.QueryRow("SELECT * FROM users WHERE (name = $1 AND password = $2)", user.Name, user.Password)

		if err := row.Scan(&foundUser.ID, &foundUser.Name, &foundUser.Password, &foundUser.Token, &foundUser.Refresh_token, &foundUser.User_type); err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "password or username incorrect error " + user.Name + user.Password})
				return
			}

			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		token, refreshToken, err := GenerateAllTokens(foundUser.Name, foundUser.User_type)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		UpdateAllTokens(token, refreshToken, foundUser.ID)

		newrow := db.DB.QueryRow("SELECT * FROM users WHERE (name = $1 AND password = $2)", user.Name, user.Password)

		if err := newrow.Scan(&foundUser.ID, &foundUser.Name, &foundUser.Password, &foundUser.Token, &foundUser.Refresh_token, &foundUser.User_type); err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "password or username incorrect error2"})
				return
			}
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)

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
