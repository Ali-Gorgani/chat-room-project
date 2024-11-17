package hash

import (
	"fmt"

	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/logger"
	"golang.org/x/crypto/bcrypt"
)

// HashedPassword hashes the given password and returns the hashed password or logs a fatal error using zerolog.
func HashedPassword(password string, logger *logger.Logger) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to hash password: %v", err))
		return "", err
	}
	return string(hashedPassword), err
}

// ComparePassword compares the given hashed password with the plain password and returns true if they match or logs a fatal error using zerolog.
func ComparePassword(hashedPassword, password string, logger *logger.Logger) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to compare password: %v", err))
		return false, err
	}
	return true, nil
}
