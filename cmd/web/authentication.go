package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func IsAuthenticated(token string, secretKey []byte) bool {
	_, err := VerifyToken(token, secretKey)
	return err == nil
}

func IsRole(token, roleRequired string, secretKey []byte) bool {
	role, err := VerifyToken(token, secretKey)
	return err == nil && role == roleRequired
}

func GenerateToken(role string, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(secretKey)
}

func VerifyToken(tokenString string, secretKey []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["role"].(string), nil
	}
	return "", errors.New("undefined token")
}
