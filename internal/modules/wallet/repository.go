package wallet

import (
	"context"
	"time"

	walletdto "twitter_clone/internal/modules/wallet/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetBalance(userID int64) (int64, *apperror.AppError)
	CreateEntry(userID int64, refUserID *int64, amount int64, txType string, currencyCode string) *apperror.AppError
	GetLedgerEntries(userID int64, limit, offset int) ([]walletdto.LedgerModel, *apperror.AppError)
}

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db}
}

// GetBalance returns current balance
func (r *WalletRepository) GetBalance(userID int64) (int64, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var balance int64
	err := r.db.QueryRow(ctx, `SELECT balance FROM wallets WHERE user_id=$1`, userID).Scan(&balance)
	if err != nil {
		return 0, apperror.DB("failed to get balance", err)
	}
	return balance, nil
}

func (r *WalletRepository) CreateEntry(
	userID int64,
	refUserID *int64,
	amount int64,
	txType string,
	currencyCode string,
) *apperror.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.DB("failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// بررسی وجود ارز
	const checkCurrency = `SELECT EXISTS(SELECT 1 FROM currencies WHERE code = $1)`
	var exists bool
	if err := tx.QueryRow(ctx, checkCurrency, currencyCode).Scan(&exists); err != nil {
		return apperror.DB("failed to check currency", err)
	}
	if !exists {
		return apperror.Validation("invalid currency code", nil, nil)
	}

	// ثبت رکورد در دفتر کل (Ledger)
	const insertLedger = `
		INSERT INTO ledger_entries (
			user_id,
			related_user_id,
			currency_code,
			entry_type,
			amount,
			type
		)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err = tx.Exec(
		ctx,
		insertLedger,
		userID,       // $1
		refUserID,    // $2
		currencyCode, // $3
		txType,       // $4 entry_type
		amount,       // $5
		"gift_in",    // $6 نوع کلی تراکنش
	)
	if err != nil {
		return apperror.DB("failed to insert ledger entry", err)
	}

	// به‌روزرسانی یا ساخت کیف پول
	const updateWallet = `
		INSERT INTO wallets (user_id, currency_code, balance)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, currency_code)
		DO UPDATE SET balance = wallets.balance + $3, updated_at = NOW();
	`
	_, err = tx.Exec(ctx, updateWallet, userID, currencyCode, amount)
	if err != nil {
		return apperror.DB("failed to update wallet", err)
	}

	// نهایی کردن تراکنش
	if err := tx.Commit(ctx); err != nil {
		return apperror.DB("failed to commit transaction", err)
	}

	return nil
}

// GetLedgerEntries returns paginated ledger entries
func (r *WalletRepository) GetLedgerEntries(userID int64, limit, offset int) ([]walletdto.LedgerModel, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, related_user_id, amount, entry_type, meta, created_at
		FROM ledger_entries
		WHERE user_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, apperror.DB("failed to query ledger entries", err)
	}
	defer rows.Close()

	var entries []walletdto.LedgerModel
	for rows.Next() {
		var e walletdto.LedgerModel
		err := rows.Scan(&e.ID, &e.UserID, &e.RefUserID, &e.Amount, &e.EntryType, &e.Meta, &e.CreatedAt)
		if err != nil {
			return nil, apperror.DB("failed to scan ledger entry", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
