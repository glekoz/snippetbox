package jwtAuth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("KJsPBp58VA3u9ZpQo8nfoAVin7E3c/fV9sm4ygwaJUI=")

type Sub struct {
	ID    int
	Name  string
	Email string
}

func CreateJWTToken(name, email string, id int) (string, error) {
	user := Sub{
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

func VerifyJWTToken(tokenString string) (*Sub, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		return secretKey, nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	subMap, ok := claims["sub"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	user := &Sub{
		ID:    int(subMap["ID"].(float64)),
		Name:  subMap["Name"].(string),
		Email: subMap["Email"].(string),
	}

	return user, nil
}
