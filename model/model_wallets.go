package model

import "time"

type Wallet struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	GroupID        *string   `json:"group_id"`
	WalletName     string    `json:"wallet_name"`
	AccountNumber  *string   `json:"account_number"`
	InitialBalance float64   `json:"initial_balance"`
	CurrentBalance float64   `json:"current_balance"`
	IsShared       bool      `json:"is_shared"`
	CreatedAt      time.Time `json:"created_at"`
}
