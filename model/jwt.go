package model

import (
	"fmt"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const JWT_ISS = "server"

func IssueJWToken(uid uint, roles []string, exp time.Time) (string, error) {
	if len(roles) == 0 {
		roles = []string{"user"}
	}

	claims := jwt.MapClaims{}
	claims["iss"] = JWT_ISS
	claims["user"] = uid
	claims["roles"] = roles
	claims["exp"] = exp.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("JWT_KEY_SIGNIN"))
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return tokenString, fmt.Errorf("Couldn't issue token: %v", err)
	}

	return tokenString, nil
}

func IssueJWTTokenForEmail(uid uint, email string, exp time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := jwt.MapClaims{}
	claims["iss"] = JWT_ISS
	claims["user"] = uid
	claims["email"] = email
	claims["exp"] = exp.Unix()
	token.Claims = claims
	secretKey := []byte(os.Getenv("JWT_KEY_EMAIL"))
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return tokenString, fmt.Errorf("Couldn't issue token: %v", err)
	}

	return tokenString, nil
}

func JWTTokenExpirationDate() time.Time {
	var err error
	var jwtHours int
	jwtHours, err = strconv.Atoi(os.Getenv("JWT_TOKEN_EXPIRATION"))
	if err != nil {
		jwtHours = 24
	}
	jwtDuration := time.Duration(jwtHours) * time.Hour
	return (time.Now().Add(jwtDuration)).Round(time.Millisecond)
}

func VerifyJWTToken(tokenString, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	// Print an error if we can't parse for some reason
	if err != nil {
		return token, fmt.Errorf("Couldn't parse token: %v", err)
	}

	// Is token invalid?
	if !token.Valid {
		return token, fmt.Errorf("Token is invalid")
	}

	return token, nil
}
