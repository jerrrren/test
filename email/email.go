package email

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bojie/orbital/backend/auth"
	"github.com/bojie/orbital/backend/db"
	"github.com/gin-gonic/gin"
	"github.com/go-mail/mail"
)

type EmailToken struct {
	Token string `json:"token"`
}

func verifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {

		var emailToken EmailToken

		if err := c.BindJSON(&emailToken); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims, err1 := auth.ValidateEmailToken(emailToken.Token)

		if err1 != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": err1})
			c.Abort()
			return
		}

		_, err2 := db.DB.Exec("UPDATE users SET verified=$1 where uid = $2 AND name = $3 ", true, claims.ID, claims.Name)

		if err2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err2.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Verification is successful, please return to the home page"})
	}
}

func sendVerificationMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("hello")
		id, ok := c.Request.URL.Query()["id"]
		uid, _ := strconv.Atoi(id[0])
		fmt.Println(uid)
		if !ok {
			fmt.Println("Url Param 'id' is missing")
			return
		}

		var target_address string
		var name string

		row := db.DB.QueryRow("SELECT uid,name,email FROM users WHERE uid = $1", uid)

		if err := row.Scan(&uid, &name, &target_address); err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no row"})
				return
			}

			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		token, err := auth.GenerateEmailVerificationToken(name, uid)

		if err != nil {
			c.IndentedJSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}



		m := mail.NewMessage()

		// Set E-Mail sender
		m.SetHeader("From", "chee_jer_en@s2013.sst.edu.sg")

		m.SetHeader("To", target_address)

		// Set E-Mail subject
		m.SetHeader("Subject", "Intronus email verification")

		// Set E-Mail body. You can set plain text or html with text/html
		m.SetBody("text/plain", "Please click the link below to verify your email\n"+"https://intronusfrontend.herokuapp.com/verifyemail/"+token)

		// Settings for SMTP server
		d := mail.NewDialer("smtp.gmail.com", 587, "intronusorbital@gmail.com", "icppswsxfkqtrfle")

		// This is only needed when SSL/TLS certificate is not valid on server.
		// In production this should be set to false.
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// Now send E-Mail
		if err := d.DialAndSend(m); err != nil {
			fmt.Println(err)
		}
	}
}

func getVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Param("user_id");
		row := db.DB.QueryRow("SELECT verified FROM users WHERE uid = $1", uid)
		var verified bool
		if err := row.Scan(&verified);err != nil{
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no row"})
				return
			}
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"verified": verified,
		})
	}
}
