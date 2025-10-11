package walletdto

import "time"

//
// Internal Structs (DB Models)
//

// LedgerModel is the raw struct mapped from DB rows
type LedgerModel struct {
	ID        int64
	UserID    int64
	RefUserID *int64
	Amount    int64
	EntryType string
	Meta      map[string]any
	CreatedAt time.Time
}

//
// Common DTOs (API Responses)
//

// BalanceRes represents the response for wallet balance
type BalanceRes struct {
	Balance int64 `json:"balance" example:"1000"`
}

// LedgerEntry represents a single transaction entry for API response
type LedgerEntry struct {
	ID        int64     `json:"id" example:"1"`
	Amount    int64     `json:"amount" example:"-50"`
	EntryType string    `json:"type" example:"gift_in"`
	RefUserID *int64    `json:"related_user_id,omitempty" example:"123"`
	CreatedAt time.Time `json:"created_at" example:"2025-10-01T12:00:00Z"`
}

// LedgerRes represents a paginated list of ledger entries
type LedgerRes struct {
	Total   int64         `json:"total" example:"2"`
	Entries []LedgerEntry `json:"entries"`
}

//
// Reward DTOs
//

// RewardReq is the payload for rewarding a user
type RewardReq struct {
	Amount       int64  `json:"amount" example:"100"`
	CurrencyCode string `json:"currency_code" example:"COIN"`
}

// RewardRes is the response after rewarding a user
type RewardRes struct {
	Message string `json:"message" example:"Reward added successfully"`
}

//
// Transfer DTOs
//

// TransferReq is the payload for transferring coins between users
type TransferReq struct {
	ToUserID     int64  `json:"to_user_id" example:"123"`
	Amount       int64  `json:"amount" example:"100000"`
	CurrencyCode string `json:"currency_code" example:"COIN"` // الزامی
}

// TransferRes — پاسخ موفق انتقال
type TransferRes struct {
	Message string `json:"message" example:"Transfer successful"`
}
