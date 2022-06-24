package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bojie/orbital/backend/db"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Action    string `json:"action"`
	Message   string `json:"message"`
	Target    string `json:"target"`
	SenderId  string `json:"senderId"`
	TimeStamp string `json:"timeStamp"`
	//	Sender  *Client
}

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const SendPrivateMessage = "send-private-message"

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}

	return json
}

func GetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {

		var messages []Message
		uid := c.Param("user_id")

		rows, err := db.DB.Query("SELECT user_id_1,user_id_2,body,messagetime FROM chats WHERE user_id_1 = $1 OR user_id_2 = $1", uid)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"GetMessages1": err})
			return
		}

		for rows.Next() {
			var message Message
			if err := rows.Scan(&message.SenderId, &message.Target, &message.Message, &message.TimeStamp); err != nil {
				c.IndentedJSON(http.StatusNotFound, gin.H{"GetMessages2": err})
				return
			}

			messages = append(messages, message)
		}
		c.IndentedJSON(http.StatusOK, messages)
		defer rows.Close()

	}
}
