package model

import "time"

type History struct {
	ID            string     `json:"id"`
	WalletID      string     `json:"wallet_id"`
	TransactionBy string     `json:"transaction_by"`
	Status        string     `json:"status"`
	TransactionAt *time.Time `json:"transaction_at"`
	Amount        int64      `json:"amount"`
	ReferenceID   string     `json:"reference_id"`
	Type          string     `json:"type"`
	CreatedAt     *time.Time `json:"created_at"`
}

type Deposit struct {
	ID          string     `json:"id"`
	Status      string     `json:"status"`
	DepositedBy string     `json:"deposited_by"`
	DepositedAt *time.Time `json:"deposited_at"`
	Amount      int64      `json:"amount"`
	ReferenceId string     `json:"reference_id"`
}

type Withdrawal struct {
	ID          string     `json:"id"`
	Status      string     `json:"status"`
	WithdrawnBy string     `json:"withdrawn_by"`
	WithdrawnAt *time.Time `json:"withdrawn_at"`
	Amount      int64      `json:"amount"`
	ReferenceId string     `json:"reference_id"`
}

// TransactionParams  used for body request
type TransactionParams struct {
	WalletID      string     `json:"wallet_id"`
	TransactionBy string     `json:"transaction_by"`
	Status        string     `json:"status"`
	TransactionAt *time.Time `json:"transaction_at"`
	Amount        int64      `json:"amount"`
	ReferenceID   string     `json:"reference_id"`
	Type          string     `json:"type"`
}
