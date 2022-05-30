package main

import (
	"fmt"
	"time"
	"os"
	

	jwt "github.com/dgrijalva/jwt-go"
	
)

type SignedDetails struct {
	Name     string `json:"username"`
	User_type string `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	jwt.StandardClaims	
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(name string,user_type string)(signedToken string,signedRefreshToken string,err error){
	claims := &SignedDetails{
		Name : name,
		User_type : user_type,
		StandardClaims:jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},

	}
	refreshClaims := &SignedDetails{
		StandardClaims:jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},

	}
	token,err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		panic(err) 
	}

	return token, refreshToken, err

}

func ValidateToken(signedToken string) (claims *SignedDetails,msg string){
	token,err := jwt.ParseWithClaims( 
		signedToken,
		&SignedDetails{},	
		func(token *jwt.Token)(interface{},error){
			return []byte(SECRET_KEY),nil
		},	
	)
	if err!=nil{
		msg = err.Error()
		return
	}
	claims,ok := token.Claims.(*SignedDetails)
	if(!ok){
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}

	return claims,msg
}

func UpdateAllTokens(signedToken string,signedRefreshToken string,userId uint){
	_, err := db.Exec("UPDATE users SET token = $1, refresh_token = $2 WHERE uid = $3",signedToken, signedRefreshToken,userId)
	if err != nil {
		panic(err)
	}
}