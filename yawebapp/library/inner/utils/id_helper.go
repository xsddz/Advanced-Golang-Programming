package utils

import "github.com/google/uuid"

func GenrateRequestID() string {
	return uuid.New().String()
}
