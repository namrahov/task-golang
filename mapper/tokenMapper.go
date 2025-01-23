package mapper

import (
	"github.com/google/uuid"
	"task-golang/model"
	"time"
)

func BuildActivationToken(activationToken string, userId int64) *model.Token {
	token := &model.Token{
		Token:           uuid.New().String(), // Generate a unique token using UUID
		ActivationToken: activationToken,
		UserID:          userId,
		CreatedAt:       time.Now(),
		TTL:             3600, // 1 hour
	}

	return token
}
