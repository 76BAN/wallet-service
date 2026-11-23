package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"wallets_service/model"
	"wallets_service/service"

	"github.com/labstack/echo/v5"
)

type WalletHandler struct {
	WalletService *service.WalletService
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{
		WalletService: walletService,
	}
}

func getUserID(c *echo.Context) (int64, error) {
	userID, ok := c.Get("userID").(int64)
	if !ok || userID <= 0 {
		return 0, errors.New("user ID tidak ditemukan dari token")
	}

	return userID, nil
}

func (h *WalletHandler) CreateWallet(c *echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": err.Error(),
		})
	}

	wallet := new(model.Wallet)

	if err := c.Bind(wallet); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Format request tidak valid",
		})
	}

	wallet.UserID = userID
	result, err := h.WalletService.CreateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Wallet berhasil dibuat",
		"data":    result,
	})
}

func (h *WalletHandler) FindWalletByID(c *echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": err.Error(),
		})
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID Wallet tidak valid",
		})
	}

	wallet, err := h.WalletService.FindWalletByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Wallet tidak ditemukan",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	// Untuk sementara, hanya pemilik Wallet yang dapat mengaksesnya.
	if wallet.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Kamu tidak memiliki akses ke Wallet ini",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil mengambil data Wallet",
		"data":    wallet,
	})
}

func (h *WalletHandler) FindWalletsByGroupID(c *echo.Context) error {
	_, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": err.Error(),
		})
	}

	groupID := c.Param("groupId")

	wallets, err := h.WalletService.FindWalletsByGroupID(groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil mengambil Wallet berdasarkan Group",
		"data":    wallets,
	})
}

func (h *WalletHandler) FindWalletsByUserIDAndGroupID(c *echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": err.Error(),
		})
	}

	groupID := c.Param("groupId")

	wallets, err := h.WalletService.FindWalletsByUserIDAndGroupID(
		userID,
		groupID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil mengambil Wallet pengguna",
		"data":    wallets,
	})
}

func (h *WalletHandler) FindWalletByName(c *echo.Context) error {
	_, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": err.Error(),
		})
	}

	walletName := c.QueryParam("wallet_name")

	wallet, err := h.WalletService.FindWalletByName(walletName)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Wallet berhasil ditemukan",
		"data":    wallet,
	})
}

func (h *WalletHandler) UpdateWallet(c *echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": err.Error(),
		})
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID Wallet tidak valid",
		})
	}

	existingWallet, err := h.WalletService.FindWalletByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Wallet tidak ditemukan",
		})
	}

	if existingWallet.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Kamu tidak memiliki akses untuk mengubah Wallet ini",
		})
	}

	wallet := new(model.Wallet)
	if err := c.Bind(wallet); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Format request tidak valid",
		})
	}

	wallet.ID = id
	wallet.UserID = userID

	result, err := h.WalletService.UpdateWallet(wallet)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Wallet berhasil diperbarui",
		"data":    result,
	})
}

func (h *WalletHandler) DeleteWallet(c *echo.Context) error {
	userID, err := getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"message": err.Error(),
		})
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ID Wallet tidak valid",
		})
	}

	existingWallet, err := h.WalletService.FindWalletByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Wallet tidak ditemukan",
		})
	}

	if existingWallet.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Kamu tidak memiliki akses untuk menghapus Wallet ini",
		})
	}

	err = h.WalletService.DeleteWallet(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Wallet berhasil dihapus",
	})
}
