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
	Year   	   int
	Location   string
	Faculty    string
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

func FillAndMatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get the user info
		type UserInfo struct {
			ID         int
			Name       string 
			Commitment int
			Year   	   int
			Location   string
			Faculty    string
			SameFaculty bool
		}
		user := new(UserInfo) 
		if err := c.Bind(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//insert info into singleusers table
		_, err := db.DB.Exec("INSERT INTO singleusers(name, commitment, year, location, faculty) VALUES($1, $2, $3, $4, $5) ON CONFLICT (name) DO UPDATE SET (commitment, year, location, faculty) = ($2, $3, $4, $5)", user.Name, user.Commitment, user.Year, user.Location, user.Faculty)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//try to find partner using indicators
		count := 1
		min := 100
		partner := new(Singleusers)
		rows, err := db.DB.Query("SELECT * FROM singleusers WHERE name != $1", user.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			temp := new(Singleusers)
			if err := rows.Scan(&temp.ID, &temp.Name, &temp.Commitment, &temp.Year, &temp.Location, &temp.Faculty); err != nil {
				//errror in scanning since the values are not filled, so we just continue
				continue
			}
			if (user.SameFaculty && temp.Faculty != user.Faculty) {
				continue
			}
			score := 0
			count++
			if (temp.Faculty != user.Faculty) {
				score += 3
			}
			if (temp.Location != user.Location) {
				score += 2
			}
			score += int(math.Abs(float64(temp.Commitment - user.Commitment)))
			if (score < min) {
				min = score;
				partner = temp
			}
		}

		//update singleusers and pairedusers table if user is paried/not paired
		if (count <= 2) {
			c.IndentedJSON(http.StatusOK, gin.H{
				"result":  false,
				"message": "",
			})
		} else {
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
