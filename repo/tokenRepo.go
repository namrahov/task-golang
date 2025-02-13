package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"strconv"
	"task-golang/model"
	"time"
)

type ITokenRepo interface {
	SaveToken(ctx context.Context, token *model.Token) error
	FindTokenByActivationToken(ctx context.Context, activationToken string) (*model.Token, error)
	FindTokenByUserId(ctx context.Context, userId int64) (*model.Token, error)
	DeleteToken(ctx context.Context, token *model.Token) error
	FindTokenByID(ctx context.Context, tokenID string) (*model.Token, error)
	FindTokenByToken(ctx context.Context, token string) (*model.Token, error)
	ExistByToken(ctx context.Context, token string) bool
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
	atIndexKey := fmt.Sprintf("activationTokenIndex:%s", token.ActivationToken)
	err = RedisClient.Set(ctx, atIndexKey, token.ID, time.Duration(token.TTL)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("error saving activationToken index to Redis: %w", err)
	}

	// Save userId to secondary index
	uiIndexKey := fmt.Sprintf("userIdIndex:%s", strconv.FormatInt(token.UserID, 10))
	err = RedisClient.Set(ctx, uiIndexKey, token.ID, time.Duration(token.TTL)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("error saving userId index to Redis: %w", err)
	}

	// Save token to secondary index
	tokenIndexKey := fmt.Sprintf("tokenIndex:%s", token.Token)
	err = RedisClient.Set(ctx, tokenIndexKey, token.ID, time.Duration(token.TTL)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("error saving token index to Redis: %w", err)
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

func (tr TokenRepo) FindTokenByUserId(ctx context.Context, userId int64) (*model.Token, error) {
	indexKey := fmt.Sprintf("userIdIndex:%s", strconv.FormatInt(userId, 10))
	tokenID, err := RedisClient.Get(ctx, indexKey).Result()
	if err == redis.Nil {
		fmt.Println("No token ID found for user")
		return nil, nil // Token not found
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving token ID from Redis: %w", err)
	}

	if tokenID == "" {
		return nil, fmt.Errorf("empty token ID for user ID: %d", userId)
	}

	key := fmt.Sprintf("tokens:%s", tokenID)
	data, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("No token data found for token ID")
		return nil, nil // Token not found
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving token from Redis: %w", err)
	}

	if data == "" {
		return nil, fmt.Errorf("empty token data for token ID: %s", tokenID)
	}

	var token model.Token
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling token: %w", err)
	}

	return &token, nil
}

// DeleteToken removes a token and its associated indexes from Redis
func (tr TokenRepo) DeleteToken(ctx context.Context, token *model.Token) error {
	if token == nil {
		return fmt.Errorf("token cannot be nil")
	}

	// Construct Redis keys
	primaryKey := fmt.Sprintf("tokens:%s", token.ID)
	atIndexKey := fmt.Sprintf("activationTokenIndex:%s", token.ActivationToken)
	uiIndexKey := fmt.Sprintf("userIdIndex:%d", token.UserID)
	tokenIndexKey := fmt.Sprintf("tokenIndex:%s", token.Token)

	// Start a Redis transaction to delete all keys atomically
	pipe := RedisClient.TxPipeline()

	// Delete the primary token key
	pipe.Del(ctx, primaryKey)

	// Delete secondary indexes
	if token.ActivationToken != "" {
		pipe.Del(ctx, atIndexKey)
	}
	pipe.Del(ctx, uiIndexKey)
	if token.Token != "" {
		pipe.Del(ctx, tokenIndexKey)
	}

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error deleting token and indexes from Redis: %w", err)
	}

	return nil
}

// FindTokenByID retrieves a token from Redis using its ID
func (tr TokenRepo) FindTokenByID(ctx context.Context, tokenID string) (*model.Token, error) {
	// Construct the Redis key for the token
	key := fmt.Sprintf("tokens:%s", tokenID)

	// Retrieve the token data from Redis
	data, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Token not found
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving token from Redis: %w", err)
	}

	// Unmarshal the token data into a Token struct
	var token model.Token
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling token: %w", err)
	}

	return &token, nil
}

func (tr TokenRepo) FindTokenByToken(ctx context.Context, token string) (*model.Token, error) {
	// Construct the token index key
	tokenIndexKey := fmt.Sprintf("tokenIndex:%s", token)

	// Get the token ID from the token index
	tokenID, err := RedisClient.Get(ctx, tokenIndexKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("token not found: %s", token)
		}
		return nil, fmt.Errorf("error fetching token ID from Redis: %w", err)
	}

	// Construct the primary token key using the retrieved token ID
	key := fmt.Sprintf("tokens:%s", tokenID)

	// Get the token data from Redis
	tokenData, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("token not found for token ID: %s", tokenID)
		}
		return nil, fmt.Errorf("error fetching token data from Redis: %w", err)
	}

	// Unmarshal the token data into a Token struct
	var tokenModel model.Token
	err = json.Unmarshal([]byte(tokenData), &tokenModel)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling token data: %w", err)
	}

	return &tokenModel, nil
}

// ExistByToken checks if a token exists in Redis
func (tr TokenRepo) ExistByToken(ctx context.Context, token string) bool {
	// Construct the token index key
	tokenIndexKey := fmt.Sprintf("tokenIndex:%s", token)
	// Check if the key exists in Redis
	exists, err := RedisClient.Exists(ctx, tokenIndexKey).Result()
	if err != nil {
		_ = fmt.Errorf("error checking if token exists in Redis: %w", err)
		return false
	}
	// Redis EXISTS command returns the number of keys that exist (0 or 1 in this case)
	fmt.Println(exists > 0)
	return exists > 0
}
