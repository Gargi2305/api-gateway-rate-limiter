package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	claims := jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte("my-secret-key"))
	if err != nil {
		panic(err)
	}

	fmt.Println(signedToken)
}
