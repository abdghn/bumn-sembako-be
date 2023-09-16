/*
 * Created on 05/04/22 09.31
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package helper

import (
	"bumn-sembako-be/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"

	"github.com/spf13/viper"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass *string) {
	bytePass := []byte(*pass)
	hPass, _ := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	*pass = string(hPass)
}

func ComparePassword(dbPass, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass)) == nil
}

// GenerateToken -> generates token
func GenerateToken(user *model.User) string {

	claims := jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 3).Unix(),
		"iat":  time.Now().Unix(),
		"data": user,
		"sub":  uint(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(viper.Get("JWT_SECRET").(string)))
	return t

}

// ValidateToken --> validate the given token
func ValidateToken(token string) (*jwt.Token, error) {

	//2nd arg function return secret key after checking if the signing method is HMAC and returned key is used by 'Parse' to decode the token)
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//nil secret key
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.Get("JWT_SECRET").(string)), nil
	})
}
