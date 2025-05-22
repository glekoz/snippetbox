package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("KJsPBp58VA3u9ZpQo8nfoAVin7E3c/fV9sm4ygwaJUI=")

type sub struct {
	ID    int
	Name  string
	Email string
}

func (app *application) createJWTToken(name, email string, id int) (string, error) {
	user := sub{
		ID:    id,
		Name:  name,
		Email: email,
	}
	payload := jwt.MapClaims{
		"sub": user,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (app *application) verifyJWTToken(tokenString string) error {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
