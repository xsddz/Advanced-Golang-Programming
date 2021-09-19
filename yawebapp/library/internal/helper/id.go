package helper

import "github.com/google/uuid"

func NextRequestID() string {
	return uuid.New().String()
}
