package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	Name     string `json:"name"`
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
}

type CustomClaims struct {
	Rollno string `json:"rollno"`
	jwt.StandardClaims
}

type LoginToken struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type RewardPayload struct {
	Rollno string `json:"rollno"`
	Coins  int64  `json:"coins,string"`
}

type TransferPayload struct {
	SenderRollno   string `json:"sender"`
	ReceiverRollno string `json:"receiver"`
	Coins          int64  `json:"coins,string"`
}
