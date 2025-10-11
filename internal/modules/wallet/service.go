package wallet

import (
	walletdto "twitter_clone/internal/modules/wallet/dto"
	"twitter_clone/internal/pkg/apperror"
)

// Service interface برای هندلر
type Service interface {
	GetBalance(userID int64) (int64, *apperror.AppError)
	GetLedger(userID int64, limit, offset int) ([]walletdto.LedgerModel, *apperror.AppError)
	Reward(userID int64, amount int64, currencyCode string) *apperror.AppError
	Transfer(fromUserID, toUserID int64, amount int64, currencyCode string) *apperror.AppError
	TopUp(userID int64, amount int64, currencyCode string) *apperror.AppError
}

type WalletService struct {
	repo Repository
}

func NewWalletService(repo Repository) *WalletService {
	return &WalletService{repo}
}

// GetBalance → مستقیم از repo
func (s *WalletService) GetBalance(userID int64) (int64, *apperror.AppError) {
	return s.repo.GetBalance(userID)
}

// GetLedger → تاریخچه تراکنش‌ها
func (s *WalletService) GetLedger(userID int64, limit, offset int) ([]walletdto.LedgerModel, *apperror.AppError) {
	return s.repo.GetLedgerEntries(userID, limit, offset)
}

// Reward → افزودن ارز برای کاربر (مثلاً انجام challenge یا دریافت gift)
func (s *WalletService) Reward(userID int64, amount int64, currencyCode string) *apperror.AppError {

	// TODO فرد خودش نمیتونه که درخواست بده و موجودیش اضافه بشه باید شرط داشته باشه

	if amount <= 0 {
		return apperror.Validation("reward amount must be positive", nil, nil)
	}
	return s.repo.CreateEntry(userID, nil, amount, "reward", currencyCode)
}

// Transfer — انتقال وجه بین کاربران
func (s *WalletService) Transfer(fromUserID, toUserID int64, amount int64, currencyCode string) *apperror.AppError {
	if fromUserID == toUserID {
		return apperror.Validation("cannot transfer to yourself", nil, nil)
	}
	if amount <= 0 {
		return apperror.Validation("amount must be positive", nil, nil)
	}

	// بررسی موجودی فرستنده
	balance, err := s.repo.GetBalance(fromUserID)
	if err != nil {
		return err
	}
	if balance < amount {
		return apperror.Validation("insufficient balance", nil, nil)
	}

	// کم کردن از فرستنده
	if err := s.repo.CreateEntry(fromUserID, &toUserID, -amount, "gift_out", currencyCode); err != nil {
		return err
	}

	// اضافه کردن به گیرنده
	if err := s.repo.CreateEntry(toUserID, &fromUserID, amount, "gift_in", currencyCode); err != nil {
		return err
	}

	return nil
}

// TopUp → شارژ حساب (مثلاً خرید سکه یا واریز از TON/USDT)
func (s *WalletService) TopUp(userID int64, amount int64, currencyCode string) *apperror.AppError {
	if amount <= 0 {
		return apperror.Validation("topup amount must be positive", nil, nil)
	}
	return s.repo.CreateEntry(userID, nil, amount, "TOPUP", currencyCode)
}
