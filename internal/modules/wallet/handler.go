package wallet

import (
	"net/http"
	"strconv"

	walletdto "twitter_clone/internal/modules/wallet/dto"
	"twitter_clone/internal/pkg/apperror"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewWalletHandler(service Service) *Handler {
	return &Handler{service: service}
}

//
// GetBalance
//

// GetBalance godoc
// @Summary      Get wallet balance
// @Description  Returns the current wallet balance of the authenticated user
// @Tags         Wallet
// @Produce      json
// @Success      200 {object} walletdto.BalanceRes "Wallet balance in smallest unit"
// @Failure      401 {object} apperror.AppError "Unauthorized, missing or invalid token"
// @Failure      500 {object} apperror.AppError "Internal server error"
// @Router       /wallet/balance [get]
func (h *Handler) GetBalance(c echo.Context) error {
	userID := c.Get("userID").(int64)

	balance, err := h.service.GetBalance(userID)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	res := walletdto.BalanceRes{Balance: balance}
	return c.JSON(http.StatusOK, res)
}

//
// GetLedger
//

func mapLedgerModelsToDTO(models []walletdto.LedgerModel) []walletdto.LedgerEntry {
	dtos := make([]walletdto.LedgerEntry, len(models))
	for i, m := range models {
		dtos[i] = walletdto.LedgerEntry{
			ID:        m.ID,
			Amount:    m.Amount,
			EntryType: m.EntryType,
			RefUserID: m.RefUserID,
			CreatedAt: m.CreatedAt,
		}
	}
	return dtos
}

// GetLedger godoc
// @Summary      Get wallet transaction history
// @Description  Returns paginated transaction history (ledger entries) of the authenticated user
// @Tags         Wallet
// @Produce      json
// @Param        limit   query int false "Number of records to return (default 20)"
// @Param        offset  query int false "Offset for pagination (default 0)"
// @Success      200 {object} walletdto.LedgerRes "List of ledger entries"
// @Failure      401 {object} apperror.AppError "Unauthorized, missing or invalid token"
// @Failure      500 {object} apperror.AppError "Internal server error"
// @Router       /wallet/ledger [get]
func (h *Handler) GetLedger(c echo.Context) error {
	userID := c.Get("userID").(int64)

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	entries, err := h.service.GetLedger(userID, limit, offset)
	if err != nil {
		return c.JSON(err.StatusCode, err)
	}

	entriesDTO := mapLedgerModelsToDTO(entries)

	res := walletdto.LedgerRes{
		Total:   int64(len(entriesDTO)), // یا total واقعی از DB
		Entries: entriesDTO,
	}

	return c.JSON(http.StatusOK, res)
}

//
// Reward
//

// Reward godoc
// @Summary      Reward user with coins
// @Description  Adds coins to the authenticated user's wallet (e.g. challenge reward)
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        data body walletdto.RewardReq true "Reward payload"
// @Success      200 {object} walletdto.RewardRes "Reward added successfully"
// @Failure      400 {object} apperror.AppError "Validation failed (e.g. negative amount)"
// @Failure      401 {object} apperror.AppError "Unauthorized, missing or invalid token"
// @Failure      500 {object} apperror.AppError "Internal server error"
// @Router       /wallet/reward [post]
func (h *Handler) Reward(c echo.Context) error {
	userID := c.Get("userID").(int64)

	var req walletdto.RewardReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid payload", nil, err))
	}

	appErr := h.service.Reward(userID, req.Amount, req.CurrencyCode)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	res := walletdto.RewardRes{Message: "Reward added successfully"}
	return c.JSON(http.StatusOK, res)
}

//
// Transfer
//

// Transfer godoc
// @Summary      Transfer coins to another user
// @Description  Transfers coins from authenticated user to target user
// @Tags         Wallet
// @Accept       json
// @Produce      json
// @Param        data body walletdto.TransferReq true "Transfer payload"
// @Success      200 {object} walletdto.TransferRes "Transfer successful"
// @Failure      400 {object} apperror.AppError "Validation failed"
// @Failure      401 {object} apperror.AppError "Unauthorized"
// @Failure      500 {object} apperror.AppError "Internal server error"
// @Router       /wallet/transfer [post]
func (h *Handler) Transfer(c echo.Context) error {
	fromUserID := c.Get("userID").(int64)

	var req walletdto.TransferReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, apperror.Validation("invalid payload", nil, err))
	}

	appErr := h.service.Transfer(fromUserID, req.ToUserID, req.Amount, req.CurrencyCode)
	if appErr != nil {
		return c.JSON(appErr.StatusCode, appErr)
	}

	return c.JSON(http.StatusOK, walletdto.TransferRes{
		Message: "Transfer successful",
	})
}
