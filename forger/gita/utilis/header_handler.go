package utilis

import (
	"errors"
	"regexp"
	"strings"
)

func HeaderHandler(headers map[string]string) (string, error) {

	authHeader := headers["Authorization"]
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	email := tokenParts[1]

	// Validate email format
	if !isValidEmail(email) {
		return "", errors.New("invalid email format")
	}

	return email, nil
}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
