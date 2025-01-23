package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"task-golang/model"
	"time"
)

type ITokenRepo interface {
	SaveToken(ctx context.Context, token *model.Token) error
	FindTokenByActivationToken(ctx context.Context, activationToken string) (*model.Token, error)
}

type TokenRepo struct {
}

// SaveToken saves a token to Redis with TTL
func (tr TokenRepo) SaveToken(ctx context.Context, token *model.Token) error {
	// Auto-generate ID if not provided
	if token.ID == "" {
		token.ID = uuid.New().String() // Generate a unique UUID
	}

	// Primary token key
	key := fmt.Sprintf("tokens:%s", token.ID)
	tokenData, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("error marshalling token: %w", err)
	}

	// Save token to Redis with TTL
	err = RedisClient.Set(ctx, key, tokenData, time.Duration(token.TTL)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("error saving token to Redis: %w", err)
	}

	// Save activationToken to secondary index
	indexKey := fmt.Sprintf("activationTokenIndex:%s", token.ActivationToken)
	err = RedisClient.Set(ctx, indexKey, token.ID, time.Duration(token.TTL)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("error saving activationToken index to Redis: %w", err)
	}

	return nil
}

// FindTokenByActivationToken retrieves a token from Redis using activationToken
func (tr TokenRepo) FindTokenByActivationToken(ctx context.Context, activationToken string) (*model.Token, error) {
	// Get the token ID from the secondary index
	indexKey := fmt.Sprintf("activationTokenIndex:%s", activationToken)
	tokenID, err := RedisClient.Get(ctx, indexKey).Result()
	if err == redis.Nil {
		return nil, nil // Token not found
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving token ID from Redis: %w", err)
	}

	// Get the token data using the token ID
	key := fmt.Sprintf("tokens:%s", tokenID)
	data, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Token not found
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving token from Redis: %w", err)
	}

	var token model.Token
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling token: %w", err)
	}

	return &token, nil
}
