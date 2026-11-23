package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"wallets_service/model"
	"wallets_service/repository"
)

type WalletService struct {
	WalletsRepository *repository.WalletsRepository
}

func NewWalletService(
	walletsRepository *repository.WalletsRepository,
) *WalletService {
	return &WalletService{
		WalletsRepository: walletsRepository,
	}
}

// CreateWallet digunakan untuk membuat Wallet baru.
func (s *WalletService) CreateWallet(
	wallet *model.Wallet,
) (*model.Wallet, error) {

	if wallet == nil {
		return nil, errors.New("data Wallet tidak boleh kosong")
	}

	if wallet.UserID <= 0 {
		return nil, errors.New("user ID tidak valid")
	}

	if wallet.WalletName == "" {
		return nil, errors.New("nama Wallet wajib diisi")
	}

	if wallet.InitialBalance < 0 {
		return nil, errors.New("saldo awal tidak boleh negatif")
	}

	// Saldo saat ini pertama kali dibuat sama dengan saldo awal.
	wallet.CurrentBalance = wallet.InitialBalance

	result, err := s.WalletsRepository.CreateWallets(wallet)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat Wallet: %w", err)
	}

	return result, nil
}

// FindWalletByID digunakan untuk mencari Wallet berdasarkan ID.
func (s *WalletService) FindWalletByID(
	id int64,
) (*model.Wallet, error) {

	if id <= 0 {
		return nil, errors.New("ID Wallet tidak valid")
	}

	wallet, err := s.WalletsRepository.FindWalletsByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Wallet tidak ditemukan")
		}

		return nil, fmt.Errorf("gagal mencari Wallet: %w", err)
	}

	return wallet, nil
}

// FindWalletsByGroupID digunakan untuk mengambil semua Wallet dalam Group.
func (s *WalletService) FindWalletsByGroupID(
	groupID string,
) ([]model.Wallet, error) {

	if strings.TrimSpace(groupID) == "" {
		return nil, errors.New("Group ID tidak boleh kosong")
	}

	wallets, err := s.WalletsRepository.FindwalletsByGroupId(groupID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil Wallet berdasarkan Group: %w", err)
	}

	return wallets, nil
}

// FindWalletsByUserIDAndGroupID digunakan untuk mencari Wallet
// berdasarkan User ID dan Group ID.
func (s *WalletService) FindWalletsByUserIDAndGroupID(
	userID int64,
	groupID string,
) ([]model.Wallet, error) {

	if userID <= 0 {
		return nil, errors.New("User ID tidak valid")
	}

	if strings.TrimSpace(groupID) == "" {
		return nil, errors.New("Group ID tidak boleh kosong")
	}

	wallets, err := s.WalletsRepository.FindByuserIdAndGrupId(
		userID,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal mencari Wallet: %w", err)
	}

	return wallets, nil
}

// FindWalletByName digunakan untuk mencari Wallet berdasarkan nama.
func (s *WalletService) FindWalletByName(
	walletName string,
) (*model.Wallet, error) {

	if strings.TrimSpace(walletName) == "" {
		return nil, errors.New("nama Wallet tidak boleh kosong")
	}

	wallet, err := s.WalletsRepository.FindByWalletsName(walletName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Wallet tidak ditemukan")
		}

		return nil, fmt.Errorf("gagal mencari Wallet: %w", err)
	}

	return wallet, nil
}

// UpdateWallet digunakan untuk mengubah data Wallet.
func (s *WalletService) UpdateWallet(
	wallet *model.Wallet,
) (*model.Wallet, error) {

	if wallet == nil {
		return nil, errors.New("data Wallet tidak boleh kosong")
	}

	if wallet.ID <= 0 {
		return nil, errors.New("ID Wallet tidak valid")
	}

	if strings.TrimSpace(wallet.WalletName) == "" {
		return nil, errors.New("nama Wallet wajib diisi")
	}

	result, err := s.WalletsRepository.UpdateWallets(wallet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Wallet tidak ditemukan")
		}

		return nil, fmt.Errorf("gagal mengubah Wallet: %w", err)
	}

	return result, nil
}

// DeleteWallet digunakan untuk menghapus Wallet.
func (s *WalletService) DeleteWallet(
	id int64,
) error {

	if id <= 0 {
		return errors.New("ID Wallet tidak valid")
	}

	err := s.WalletsRepository.DeleteWallets(id)
	if err != nil {
		return fmt.Errorf("gagal menghapus Wallet: %w", err)
	}

	return nil
}
