package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type UserClaims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
}

func CreateToken(id, secret string, exp time.Duration) (string, error) {
	claims := UserClaims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.NewTime(float64(time.Now().Add(exp).Unix())),
			Issuer:    "typenotes-go-graphql",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func CreateTokens(id string) (string, string, error) {
	token, err := CreateToken(id, os.Getenv("TOKEN_SECRET"), time.Minute)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := CreateToken(id, os.Getenv("REFRESH_TOKEN"), time.Hour*24*7)
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := ParseToken(refreshToken, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		return "", "", err
	}
	return CreateTokens(claims.UserID)
}
