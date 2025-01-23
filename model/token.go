package model

import "time"

type Token struct {
	ID              string    `json:"id"`              // Redis ID
	Token           string    `json:"token"`           // token
	ActivationToken string    `json:"activationToken"` // activation token
	UserID          int64     `json:"userId"`          // Indexed user ID
	CreatedAt       time.Time `json:"createdAt"`       // Creation timestamp
	TTL             int64     `json:"ttl"`             // Time-to-live in seconds
}
