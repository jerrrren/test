package main

type User struct {
	ID            uint   `json:"uid"`
	Name          string `json:"username"`
	Password      string `json:"password"`
	User_type     string `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Refresh_token string `json:"refresh_token"`
	Token         string `json:"token"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}
