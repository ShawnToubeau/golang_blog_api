package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Creates an authorization token for a specific user ID.
func CreateToken(userId uint32) (string, error) {
	// create claims object
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	// set expiration time to 1 hour
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	// create auth token using HS256 singing method and claims object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// sign the token
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// Extract an auth token from a request header.
func ExtractToken(r *http.Request) string {
	// request query parameters
	keys := r.URL.Query()
	// todo: not sure what this is
	token := keys.Get("token")
	fmt.Printf("ExtractToken - token: %v\n", token)
	if token != "" {
		return token
	}
	// extract bearer token
	bearerToken := r.Header.Get("Authorization")
	fmt.Printf("ExtractToken - bearerToken: %v\n", token)
	if len(strings.Split(bearerToken, "")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

// Formats and prints the token claims in the console.
func FormatClaims(claim interface{}) {
	// formats claims object
	b, err := json.MarshalIndent(claim, "", " ")
	if err != nil {
		return
	}

	fmt.Println(string(b))
}

// Matches the signing method of the auth token.
func CheckSigningMethod(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
}

// Confirms whether the token is valid.
func TokenValid(r *http.Request) error {
	// extract auth token from request
	tokenString := ExtractToken(r)
	// checks signing method
	token, err := CheckSigningMethod(tokenString)
	if err != nil {
		return err
	}
	// return token claims if valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		FormatClaims(claims)
	}
	return nil
}

// Extracts the user's ID from the auth token.
func ExtractTokenID(r *http.Request) (uint32, error) {
	// get bearer token value from request header
	tokenString := ExtractToken(r)
	// parse the token string to make sure its valid
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// todo: not sure what this does
		stuff := []byte(os.Getenv("API_SECRET"))
		fmt.Printf("ExtractTokenID - stuff: %v\n", stuff)
		return stuff, nil
	})
	// throw parsing error
	if err != nil {
		return 0, err
	}
	// token claims (authorized: boolean, userID: uint32, expiration: time)
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// get user ID from claims object, converting from string to uint64
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		// return as uint32
		return uint32(uid), nil
	}
	return 0, nil
}
