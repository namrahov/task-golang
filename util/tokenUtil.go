package util

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"task-golang/model"
	"task-golang/repo"
	"time"
)

type ITokenUtil interface {
	GenerateToken() string
	GenerateSHA(text string) string
	ReSetActivationToken(ctx context.Context, user *model.User, activationToken string)
}

type TokenUtil struct {
	TokenRepo repo.ITokenRepo
}

// GenerateToken generates a random token by hashing a UUID.
func (tu *TokenUtil) GenerateToken() string {
	randomUUID := uuid.New()
	return tu.GenerateSHA(randomUUID.String())
}

// GenerateSHA generates a SHA-256 hash for the given text.
func (tu *TokenUtil) GenerateSHA(text string) string {
	hash := sha256.New()
	_, err := hash.Write([]byte(text))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (tu *TokenUtil) ReSetActivationToken(ctx context.Context, user *model.User, activationToken string) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.Register.start")

	existingToken, errExistingToken := tu.TokenRepo.FindTokenByUserId(ctx, user.Id)
	if errExistingToken != nil {
		fmt.Println("reSetActivationToken error=" + errExistingToken.Error())
	}

	if existingToken != nil {
		existingToken.ActivationToken = activationToken
		tu.TokenRepo.SaveToken(ctx, existingToken)
	} else {
		tokenEntity := &model.Token{
			ActivationToken: activationToken,
			UserID:          user.Id,
			CreatedAt:       time.Now(),
			TTL:             600, // 10 dakika
		}
		tu.TokenRepo.SaveToken(ctx, tokenEntity)
	}
}
