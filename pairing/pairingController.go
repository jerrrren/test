package pairing

import (
	"database/sql"
	"math"
	"net/http"

	"github.com/bojie/orbital/backend/db"
	"github.com/gin-gonic/gin"
)

type Singleusers struct {
	ID         int
	Name       string `json:"Name"`
	Commitment int
	Location   string
	Filledinfo bool
}

type Paireduser struct {
	ID      int
	Name    string
	Partner string
}

func checkIfPaired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := new(Paireduser)
		result := true
		
		if err := c.Bind(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in binding name to user"})
			return
		}
		row := db.DB.QueryRow("SELECT * from pairedusers WHERE name=$1", user.Name)
		if err := row.Scan(&user.ID, &user.Name, &user.Partner); err != nil {
			if err == sql.ErrNoRows {
				result = false
			} else {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			}
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": result,"partner":user.Partner})
	}
}

//function to add user to singleusers table during registration
func AddRegistrationUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := new(Singleusers)
		if err := c.Bind(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in binding name to user"})
			return
		}

		_, err := db.DB.Exec("INSERT INTO singleusers (name) VALUES ($1)", user.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error inserting (user not unique)"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "successful"})
	}
}

func FillAndMatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//update the information about the users
		user := new(Singleusers)
		if err := c.Bind(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.DB.Exec("UPDATE singleusers SET commitment=$1, location=$2, filledinfo=true WHERE name=$3", user.Commitment, user.Location, user.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//try to match user to another user with the min value (closest in terms of indicators)
		min := 20
		count := 1
		partner := new(Singleusers)
		rows, err := db.DB.Query("SELECT * FROM singleusers WHERE name != $1", user.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			temp := new(Singleusers)
			if err := rows.Scan(&temp.ID, &temp.Name, &temp.Commitment, &temp.Location, &temp.Filledinfo); err != nil {
				//errror in scanning since the values are not filled, so we just continue
				continue
			}
			score := 0
			count++
			if temp.Location != user.Location {
				score++
			}
			score += int(math.Abs(float64(temp.Commitment - user.Commitment)))
			if score < min {
				min = score
				partner = temp
			}
		}

		//if the num of users we encounter is less than 2, we return false
		if count <= 2 {
			c.IndentedJSON(http.StatusOK, gin.H{
				"result":  false,
				"message": "",
			})
		} else {
		//else we count the pairing as successful, and delete the pair from singleusers and insert into pairedusers
			_, err := db.DB.Exec("DELETE FROM singleusers WHERE name=$1 OR name=$2", user.Name, partner.Name)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			_, err = db.DB.Exec("INSERT INTO pairedusers (name, partner) VALUES ($1, $2), ($3, $4)", user.Name, partner.Name, partner.Name, user.Name) 
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.IndentedJSON(http.StatusOK, gin.H{
				"result":  true,
				"message": partner.Name,
			})
		}
	}
}

func DeletePairedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := new(Paireduser)
		if err := c.Bind(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		row := db.DB.QueryRow("SELECT * FROM pairedusers WHERE name=$1", user.Name)
		if err := row.Scan(&user.ID, &user.Name, &user.Partner); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, err := db.DB.Exec("DELETE FROM pairedusers WHERE name=$1", user.Name); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if _, err := db.DB.Exec("DELETE FROM pairedusers WHERE name=$1", user.Partner); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"message": "user and partner deleted"})
	}
}


// func DeleteSingleUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user := new(Singleusers)
// 		if err := c.Bind(user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		_, err := db.DB.Exec("DELETE FROM singleusers WHERE name=$1", user.Name)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.IndentedJSON(http.StatusOK, gin.H{"message": user.Name + " was successfully deleted from singleusers"})
// 	}
// }

// //add the 2 paired user to the pairedusers table
// func AddPaireduser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		paired := new(Paireduser)
// 		if err := c.Bind(paired); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		if _, err := db.DB.Exec("INSERT INTO pairedusers (name, partner) VALUES ($1, $2)", paired.Name, paired.Partner); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		if _, err := db.DB.Exec("INSERT INTO pairedusers (name, partner) VALUES ($1, $2)", paired.Partner, paired.Name); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.IndentedJSON(http.StatusOK, gin.H{"message": "paired users added"})
// 	}
// }

// //function to fill up the indicators for a user
// func FillIndicators() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user := new(Singleusers)
// 		if err := c.Bind(user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		_, err := db.DB.Exec("UPDATE singleusers SET commitment=$1, location=$2, filledinfo=true WHERE name=$3", user.Commitment, user.Location, user.Name)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		user.Filledinfo = true
// 		c.IndentedJSON(http.StatusOK, gin.H{"message": user.Filledinfo})
// 	}
// }

// func Match() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user := new(Singleusers)
// 		min := 20
// 		count := 1
// 		result := new(Singleusers)
// 		if err := c.Bind(user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		row := db.DB.QueryRow("SELECT * FROM singleusers WHERE name=$1", user.Name)
// 		if err := row.Scan(&user.ID, &user.Name, &user.Commitment, &user.Location, &user.Filledinfo); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "err.Error()"})
// 			return
// 		}

// 		//search for users to pair
// 		rows, err := db.DB.Query("SELECT * FROM singleusers WHERE name != $1", user.Name)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		for rows.Next() {
// 			temp := new(Singleusers)
// 			if err := rows.Scan(&temp.ID, &temp.Name, &temp.Commitment, &temp.Location, &temp.Filledinfo); err != nil {
// 				//errror in scanning since the values are not filled, so we just continue
// 				continue
// 			}
// 			score := 0
// 			count++
// 			if temp.Location != user.Location {
// 				score++
// 			}
// 			score += int(math.Abs(float64(temp.Commitment - user.Commitment)))
// 			if score < min {
// 				min = score
// 				result = temp
// 			}
// 		}
// 		if count <= 2 {
// 			c.IndentedJSON(http.StatusOK, gin.H{
// 				"result":  false,
// 				"message": "not enough users is available for pairing, please wait for more users to join",
// 			})
// 		} else {
// 			c.IndentedJSON(http.StatusOK, gin.H{
// 				"result":  true,
// 				"message": result.Name,
// 			})
// 		}
// 	}
// }

// //function to check if user is single and if they have filled up the indicators
// func checkInfo() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user := new(Singleusers)
// 		single := true
// 		filled := true
// 		if err := c.Bind(user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		row := db.DB.QueryRow("SELECT filledinfo FROM singleusers WHERE name=$1", user.Name)
// 		if err := row.Scan(&user.Filledinfo); err != nil {
// 			if err == sql.ErrNoRows {
// 				single = false
// 				filled = false
// 			} else {
// 				c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 				return
// 			}
// 		}
// 		filled = user.Filledinfo
// 		c.IndentedJSON(http.StatusOK, gin.H{
// 			"single": single,
// 			"filled": filled,
// 		})
// 	}
// }
