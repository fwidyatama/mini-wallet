package model

import "time"

// Wallet struct
type Wallet struct {
	ID         string     `json:"id"`
	OwnedBy    string     `json:"owned_by"`
	Status     string     `json:"status"`
	Balance    int64      `json:"balance"`
	EnabledAt  *time.Time `json:"enabled_at,omitempty"`
	DisabledAt *time.Time `json:"disabled_at,omitempty"`
}

// OwnerWalletParam is to Create wallet param struct
type OwnerWalletParam struct {
	Owner string
}

type UpdateBalanceParam struct {
	WalletId string `json:"wallet_id"`
	Balance  int64  `json:"balance"`
}
