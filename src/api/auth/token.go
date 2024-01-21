package auth

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
	Author : Mehran farhadi bajestani
	Created_at : 2024-01-20

*/

// CreateToken generates a JWT with the specified user ID and sets the expiration time.
// It returns the signed JWT string or an error.
func CreateToken(userId uint32, username, email string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{}

	// set is authenticate
	claims["authorized"] = true

	// set User id
	claims["userId"] = userId

	// set username in token
	claims["username"] = username

	//set email in token
	claims["email"] = email

	//set isAdmin in token
	claims["is_admin"] = isAdmin

	claims["exp_date"] = time.Now().Add(time.Hour * 1).Unix()
	// Token expires after 1 hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

// TokenValid checks if the provided JWT in the HTTP request is valid.
// It returns an error if the token is invalid or expired.
func TokenValid(r *http.Request) error {
	tokenStr := ExtractToken(r)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}
	return nil
}

// ExtractToken extracts the JWT from the request's URL query or Authorization header.
func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenID extracts the user ID from the JWT in the request.
// It returns the user ID or an error.
func ExtractTokenID(r *http.Request) (uint32, error) {

	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}

//func IsAdmin(r *http.Request) (bool, error) {
//	tokenString := ExtractToken(r)
//
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
//		}
//		return []byte(os.Getenv("API_SECRET")), nil
//	})
//	if err != nil {
//		return false, err
//	}
//	claims, ok := token.Claims.(jwt.MapClaims)
//	if ok && token.Valid {
//		isAdmin := claims["is_admin"]
//		return isAdmin, nil
//	}
//	return false, nil
//}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b))
}
