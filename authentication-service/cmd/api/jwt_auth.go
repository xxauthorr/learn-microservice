package main

import (
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	// "github.com/joho/godotenv"
)

type SignedCredentials struct {
	Email string
	jwt.StandardClaims
}

type Jwt struct {
	AccessToken  string `json:"access_token"`
	TokenExpiry  string `json:"expiry"`
	RefreshToken string `json:"refresh_token"`
}

type envKeys struct {
	accessKey  string
	refreshKey string
}

func loadEnv() (envKeys, error) {
	// loads env file
	var env envKeys
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	return env, err
	// }
	env.accessKey = os.Getenv("ACCESS_KEY")
	env.refreshKey = os.Getenv("REFRESH_KEY")
	return env, nil
}

func (j Jwt) GenerateToken(email string) (Jwt, error) {
	var Jwt Jwt
	key, err := loadEnv()
	if err != nil {
		log.Println("error in load env")
		return Jwt, err
	}
	exp := time.Now().Local().Add(time.Minute * time.Duration(1))
	Jwt.TokenExpiry = exp.Format("2-01-2006 3:04:05 PM")
	accessClaims := &SignedCredentials{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	// refreshClaims := &SignedCredentials{
	// 	Email: email,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(720)).Unix(),
	// 	},
	// }
	fmt.Println(key.accessKey, " ", key.refreshKey)
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(key.accessKey))
	if err != nil {
		log.Println("error in load access token")
		return Jwt, err
	}

	// refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).SignedString([]byte(key.refreshKey))
	// if err != nil {
	// 	log.Println("error in load refresh token")
	// 	return Jwt, err
	// }
	Jwt.AccessToken = accessToken
	Jwt.RefreshToken = "refreshToken"

	return Jwt, nil
}

func (j Jwt) ValidateToken(breed, token string) (claims *SignedCredentials, msg string) {
	key, err := loadEnv()
	if err != nil {
		log.Panic(err)
		return
	}
	if breed == "access token" {
		token, err := jwt.ParseWithClaims(
			token,
			&SignedCredentials{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(key.accessKey), nil
			},
		)
		if err != nil {
			msg = "	"
			return
		}
		var ok bool
		claims, ok = token.Claims.(*SignedCredentials)
		if claims.ExpiresAt < time.Now().Local().Unix() {
			msg = "Token is expired"
			return
		}
		if !ok {
			msg = "Token is invalid"
			return
		}
		return claims, msg
	}
	toks, err := jwt.ParseWithClaims(
		token,
		&SignedCredentials{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(key.refreshKey), nil
		},
	)
	if err != nil {
		msg = "error while validation"
		return
	}
	var ok bool
	claims, ok = toks.Claims.(*SignedCredentials)
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is expired"
		return
	}
	if !ok {
		msg = "Token is invalid"
		return
	}
	return claims, msg
}
