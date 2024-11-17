package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// CreateToken creates a new JWT token.
func CreateToken(secretKey string, claim UserClaims) (string, error) {
	claims, err := NewUserClaims(claim)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// VerifyToken verifies the JWT token.
func VerifyToken(tokenStr, secretKey string) (UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return UserClaims{}, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return UserClaims{}, fmt.Errorf("invalid token claims")
	}

	return *claims, nil
}
