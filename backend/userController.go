package main

import (
	"database/sql"
	"net/http"
	"fmt"


	"github.com/gin-gonic/gin"
	//"github.com/go-playground/validator/v10"

)

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

		result, err := db.Exec("INSERT INTO users (name,password,refresh_token,token,user_type) VALUES ($1, $2, $3,$4,$5)", user.Name ,user.Password,user.Refresh_token,user.Token,user.User_type)

		if err != nil {
			fmt.Printf("test")
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
		row := db.QueryRow("SELECT * FROM users WHERE (name = $1 AND password = $2)", user.Name, user.Password)

		if err := row.Scan(&foundUser.ID, &foundUser.Name, &foundUser.Password,&foundUser.Token ,&foundUser.Refresh_token, &foundUser.User_type); err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "password or username incorrect error " + user.Name+user.Password})
				return
			}

			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}

		token, refreshToken, err := GenerateAllTokens(foundUser.Name, foundUser.User_type)
		if(err!= nil){
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		
		UpdateAllTokens(token, refreshToken, foundUser.ID)

		newrow := db.QueryRow("SELECT * FROM users WHERE (name = $1 AND password = $2)", user.Name, user.Password)

		if err := newrow.Scan(&foundUser.ID, &foundUser.Name, &foundUser.Password,&foundUser.Token ,&foundUser.Refresh_token, &foundUser.User_type); err != nil {
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

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		var users []User

		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error"})
			return
		}

		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Name, &user.Password, &user.Token,&user.Refresh_token,&user.User_type); err != nil {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "error"})
				return
			}

			users = append(users, user)
		}

		c.IndentedJSON(http.StatusOK, users)
		defer rows.Close()

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Param("user_id")
		var user User

		row := db.QueryRow("SELECT * FROM users WHERE uid = $1", uid)

		if err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Token,&user.Refresh_token,&user.User_type); err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no row"})
				return
			}

			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, user)
	}
}
