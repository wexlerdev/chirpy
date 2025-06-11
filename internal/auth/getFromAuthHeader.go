package auth

import (
	"net/http"
	"strings"
	"fmt"
)

func GetBearerToken(headers http.Header) (string, error){
	headerString := headers.Get("Authorization")
	components := strings.Fields(headerString)
	if len(components) != 2 {
		return "", fmt.Errorf("Length of components is > 2: Len: %v", len(components))
	}

	if components[0] != "Bearer" {
		return "", fmt.Errorf("Prefix is not 'Bearer' Prefix: %v", components[0])
	}

	return components[1], nil

}

func GetAPIKey(headers http.Header) (string, error){
	headerString := headers.Get("Authorization")
	components := strings.Fields(headerString)
	if len(components) != 2 {
		return "", fmt.Errorf("Length of components is > 2: Len: %v", len(components))
	}

	if components[0] != "ApiKey" {
		return "", fmt.Errorf("Prefix is not 'Bearer' Prefix: %v", components[0])
	}

	return components[1], nil
}
