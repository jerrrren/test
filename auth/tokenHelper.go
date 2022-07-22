package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/bojie/orbital/backend/db"

	jwt "github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Name      string `json:"username"`
	User_type string `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	jwt.StandardClaims
}

type SignedEmailVericationDetails struct {
	Name  string `json:"username"`
	ID    int    `json:uid`
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateEmailVerificationToken(name string , uid int) (signedToken string,err error) {
	claims := &SignedEmailVericationDetails{
		Name:  name,
		ID:   uid,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(3)).Unix(),
		},
	}
	token,err := jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY))
	if err!=nil{
		fmt.Println(err)
	}

	return token,err
}

func GenerateAllTokens(name string, user_type string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Name:      name,
		User_type: user_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		panic(err)
	}

	return token, refreshToken, err

}

func ValidateEmailToken(signedToken string) (claims *SignedEmailVericationDetails,msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedEmailVericationDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}
	
	claims, ok := token.Claims.(*SignedEmailVericationDetails)
	if !ok {
		msg = "the token is invalid"
		fmt.Println("the token is invalid")
		return claims,msg
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		fmt.Println("token is expired")
		return claims,msg
	}

	return claims, msg
}



func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId uint) {
	_, err := db.DB.Exec("UPDATE users SET token = $1, refresh_token = $2 WHERE uid = $3", signedToken, signedRefreshToken, userId)
	if err != nil {
		panic(err)
	}
}
