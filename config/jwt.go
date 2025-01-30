package config

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	logger "github.com/okyws/go-banking-lib/config"
)

type Claims struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Secret key for JWT
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Generate JWT Token
func GenerateJWT(id, username string) (string, error) {
	expirationTime := time.Now().Add(2 * time.Hour)

	claims := &Claims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	logger.GetLog().Info().
		Str("username", username).
		Msg("Token generated successfully")

	return tokenString, nil
}

// ParseToken for validating JWT
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	logger.GetLog().Info().
		Str("username", claims.Username).
		Msg("Token parsed successfully")

	return claims, nil
}
