package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	value := headers.Get("Authorization")

	if value == " " {
		return "", fmt.Errorf("Cannot find authorization")
	}
	if !strings.Contains(value, "ApiKey") {
		return "", fmt.Errorf("Cannot find apikey")
	}

	apiKey := strings.TrimSpace(strings.TrimPrefix(value, "ApiKey"))
	if apiKey == " " {
		return "", fmt.Errorf("Api Key doesn't exists")
	}

	return apiKey, nil

}
