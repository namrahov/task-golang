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
	fmt.Println("islediiii")
	if token.ID == "" {
		token.ID = uuid.New().String() // Generate a unique UUID
	}

	// Primary token key
	key := fmt.Sprintf("tokens:%s", token.ID)
	tokenData, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("error marshalling token: %w", err)
	}
	fmt.Println("islediiii2")
	// Save token to Redis with TTL
	err = RedisClient.Set(ctx, key, tokenData, time.Duration(token.TTL)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("error saving token to Redis: %w", err)
	}

	// Save activationToken to secondary index
	fmt.Println("islediiii3")

	indexKey := fmt.Sprintf("activationTokenIndex:%s", token.ActivationToken)
	err = RedisClient.Set(ctx, indexKey, token.ID, time.Duration(token.TTL)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("error saving activationToken index to Redis: %w", err)
	}
	fmt.Println("islediiii4")

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

//func main() {
//	ctx := context.Background()
//
//	// Example: Save a token without manually providing an ID
//	token := &Token{
//		Token:          uuid.New().String(), // Generate a unique token using UUID
//		ActivationToken: uuid.New().String(), // Generate a unique activationToken using UUID
//		UserID:         42,
//		CreatedAt:      time.Now(),
//		TTL:            3600, // 1 hour
//	}
//
//	err := SaveToken(ctx, token)
//	if err != nil {
//		fmt.Println("Error saving token:", err)
//		return
//	}
//	fmt.Printf("Token saved successfully! Generated ID: %s\n", token.ID)
//
//	// Example: Find token by activationToken
//	retrievedToken, err := FindTokenByActivationToken(ctx, token.ActivationToken)
//	if err != nil {
//		fmt.Println("Error finding token by activationToken:", err)
//		return
//	}
//	if retrievedToken == nil {
//		fmt.Println("No token found for the given activationToken")
//		return
//	}
//	fmt.Printf("Retrieved Token: %+v\n", retrievedToken)
//}
