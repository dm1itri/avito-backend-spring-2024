package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

func (app *Application) isAuthenticated(r *http.Request) bool {
	token := r.Header.Get("token")
	_, err := verifyToken(token, app.secretKeyJWT)
	return err == nil
}

func (app *Application) isRole(r *http.Request, roleRequired string) bool {
	token := r.Header.Get("token")
	role, err := verifyToken(token, app.secretKeyJWT)
	return err == nil && role == roleRequired
}

func generateToken(role string, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	})
	return token.SignedString(secretKey)
}

func verifyToken(tokenString string, secretKey []byte) (string, error) {
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
