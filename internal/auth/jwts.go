package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("Secret length: %d\n", len(tokenSecret))
	s, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Error signing token : %v\n", err)
	}

	return s, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}

		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w\n", err)
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("Invalid token\n")
	}

	strID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error getting id")
	}
	uuidID, err := uuid.Parse(strID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid id")
	}

	return uuidID, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	value := headers.Get("Authorization")

	if value == "" {
		return "", fmt.Errorf("Cannot find authorization")
	}
	if !strings.Contains(value, "Bearer") {
		return "", fmt.Errorf("Bearer doesn't exists")
	}

	bearerKey := strings.TrimSpace(strings.TrimPrefix(value, "Bearer"))
	if bearerKey == "" {
		return "", fmt.Errorf("Bearer doesn't exists")
	}

	return bearerKey, nil
}
