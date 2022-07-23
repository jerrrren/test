package forgetpassword

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/bojie/orbital/backend/auth"
	"github.com/bojie/orbital/backend/db"
	"github.com/gin-gonic/gin"
	"github.com/go-mail/mail"
)

func sentResetEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		type Response struct {
			Username string `json:"username"`
		}
		var response Response

		if err := c.BindJSON(&response); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}


		row := db.DB.QueryRow("SELECT name,password,email FROM users WHERE name = $1", response.Username)

		var foundUser auth.User

		if err := row.Scan(&foundUser.Name, &foundUser.Password, &foundUser.Email); err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no row"})
				return
			}

			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		token, err := auth.GeneratePasswordToken (foundUser.Name,foundUser.Password)



		if err != nil {
			c.IndentedJSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}

		m := mail.NewMessage()

		// Set E-Mail sender
		m.SetHeader("From", "chee_jer_en@s2013.sst.edu.sg")

		m.SetHeader("To", foundUser.Email)

		// Set E-Mail subject
		m.SetHeader("Subject", "Intronus password reset")

		// Set E-Mail body. You can set plain text or html with text/html
		m.SetBody("text/plain", "Please click the link below to reset your password\n"+"https://intronusfrontend.herokuapp.com/updatepassword/"+token)

		// Settings for SMTP server
		d := mail.NewDialer("smtp.gmail.com", 587, "chee_jer_en@s2013.sst.edu.sg", "edplfjwgcyunfdkt")

		// This is only needed when SSL/TLS certificate is not valid on server.
		// In production this should be set to false.
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// Now send E-Mail
		if err := d.DialAndSend(m); err != nil {
			fmt.Println(err)
		}
	}
}

type ResetPasswordToken struct {
	Token string `json:"token"`
	NewPassword string `json:"password"`
}

func resetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {

		var resetPasswordToken ResetPasswordToken

		if err := c.BindJSON(&resetPasswordToken); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims, err1 := auth.ValidatePasswordVerificationToken(resetPasswordToken.Token)

		if err1 != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": err1})
			c.Abort()
			return
		}

		password := auth.HashPassword(resetPasswordToken.NewPassword)

		_, err2 := db.DB.Exec("UPDATE users SET password = $1 where name = $3 ", password, claims.Name)

		if err2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err2.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Password Reset is Successful"})
	}
}
