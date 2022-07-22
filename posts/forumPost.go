package posts

import (
	"net/http"
	"strconv"

	_ "database/sql"
	"errors"
	_ "errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	"github.com/bojie/orbital/backend/auth"
	"github.com/bojie/orbital/backend/db"
)

type UserName struct {
	ID   uint   `json:"uid"`
	Name string `json:"username"`
}
type Post struct {
	ID                   int
	Field                string `binding:"required"`
	UID                  int    `binding:"required"`
	Intro                string `binding:"required"`
	Content              string `binding:"required"`
	Name                 string
	CreatedAt            time.Time
	ModifiedAt           time.Time
	NumParticipants      int
	Participants         pq.Int32Array
	ParticipantsUsername []UserName
}

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var post Post
		if err := c.BindJSON(&post); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := db.DB.Exec("INSERT INTO posts (field,uid,intro,content,participants,num_participants) VALUES ($1, $2, $3, $4, $5,$6)", post.Field, post.UID, post.Intro, post.Content, pq.Array([]int{}),post.NumParticipants)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{"test": err})
			return
		}

		c.IndentedJSON(http.StatusOK, result)

	}

}

func getPosts() ([]Post, error) {
	var posts []Post
	rows, err := db.DB.Query("SELECT posts.id,posts.field,posts.uid,posts.intro,posts.content,posts.participants,users.name FROM posts JOIN users ON users.uid = posts.uid")
	if err != nil {
		return nil, errors.New("There is no posts in database")
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Field, &post.UID, &post.Intro, &post.Content, &post.Participants, &post.Name)
		if err != nil {
			return nil, errors.New("error is scanning and assigning values from the posts")
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPostById(id int) (*Post, error) {
	post := new(Post)
	row := db.DB.QueryRow("SELECT posts.id,posts.field,posts.uid,posts.intro,posts.content,posts.participants,users.name,posts.num_participants FROM posts JOIN users ON users.uid = posts.uid WHERE posts.id = $1", id)

	err := row.Scan(&post.ID, &post.Field, &post.UID, &post.Intro, &post.Content, &post.Participants, &post.Name,&post.NumParticipants)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("unable to fetch post: invalid post id")
	}

	var userNames []UserName

	rows, err := db.DB.Query("SELECT t1.*,users.name FROM users JOIN( SELECT UNNEST(posts.participants) FROM posts WHERE posts.id = $1) as t1 ON users.uid = t1.unnest", id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		var userName UserName
		if err := rows.Scan(&userName.ID, &userName.Name); err != nil {
			fmt.Println(err)
			return nil, err
		}

		userNames = append(userNames, userName)
	}

	post.ParticipantsUsername = userNames

	defer rows.Close()

	return post, nil
}

func GetPosts(c *gin.Context) {
	posts, err := getPosts()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func GetPostsById(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error in converting id to int"})
	}
	post, err := getPostById(intId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, post)
}

func UpdateParticipants() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		intId, err := strconv.Atoi(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error in converting id to int"})
		}

		var addedUser auth.User

		if err := c.BindJSON(&addedUser); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if(addedUser.ID == 0) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := db.DB.Exec("UPDATE posts SET participants = array_append(participants,$1) WHERE id = $2", addedUser.ID, intId)

		fmt.Println(err)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"test": err})
			return
		}

		c.IndentedJSON(http.StatusOK, result)

	}
}
