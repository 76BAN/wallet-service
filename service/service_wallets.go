package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"wallets_service/model"
	"wallets_service/rabbitmq"
	"wallets_service/repository"
)

type WalletService struct {
	WalletsRepository *repository.WalletsRepository
	LoggingPublisher  *rabbitmq.LoggingPublisher
}

func NewWalletService(
	walletsRepository *repository.WalletsRepository,
	loggingPublisher *rabbitmq.LoggingPublisher,
) *WalletService {
	return &WalletService{
		WalletsRepository: walletsRepository,
		LoggingPublisher:  loggingPublisher,
	}
}

// CreateWallet digunakan untuk membuat Wallet baru.
func (s *WalletService) CreateWallet(
	wallet *model.Wallet,
	userID int64,
) (*model.Wallet, error) {

	if wallet == nil {
		s.publishUserLog(
			"WARN",
			"CREATE_WALLET",
			"FAILED",
			"Ada Kolom Atribut yang kosong",
			&userID,
			map[string]any{
				"reason": "FAILED_CREATE_WALLET",
			},
		)
		return nil, errors.New("data Wallet tidak boleh kosong")
	}

	if wallet.UserID <= 0 {
		s.publishUserLog(
			"WARN",
			"CREATE_WALLET",
			"FAILED",
			"Ada Kolom Atribut yang kosong",
			&userID,
			map[string]any{
				"reason": "FAILED_CREATE_WALLET",
			},
		)
		return nil, errors.New("user ID tidak valid")
	}

	if wallet.WalletName == "" {
		s.publishUserLog(
			"WARN",
			"CREATE_WALLET",
			"FAILED",
			"Ada Kolom Atribut yang kosong",
			&userID,
			map[string]any{
				"reason": "FAILED_CREATE_WALLET",
			},
		)
		return nil, errors.New("nama Wallet wajib diisi")
	}

	if wallet.InitialBalance < 0 {
		s.publishUserLog(
			"WARN",
			"CREATE_WALLET",
			"FAILED",
			"Ada Kolom Atribut yang kosong",
			&userID,
			map[string]any{
				"reason": "FAILED_CREATE_WALLET",
			},
		)
		return nil, errors.New("saldo awal tidak boleh negatif")
	}

	// Saldo saat ini pertama kali dibuat sama dengan saldo awal.
	wallet.CurrentBalance = wallet.InitialBalance

	result, err := s.WalletsRepository.CreateWallets(wallet)
	if err != nil {
		s.publishUserLog(
			"WARN",
			"CREATE_WALLET",
			"FAILED",
			"Gagal Menbuat Wallet",
			&userID,
			map[string]any{
				"reason": "FAILED_CREATE_WALLET",
			},
		)
		return nil, fmt.Errorf("gagal membuat Wallet: %w", err)
	}
	s.publishUserLog(
		"INFO",
		"CREATE_WALLET",
		"SUCCESS",
		"Berhasil menbuat Wallet",
		&userID,
		map[string]any{
			"reason": "SUCCESS_CREATE_WALLET",
		},
	)

	return result, nil
}

// FindWalletByID digunakan untuk mencari Wallet berdasarkan ID.
func (s *WalletService) FindWalletByID(
	id int64,
	userID int64,
) (*model.Wallet, error) {

	if id <= 0 {
		s.publishUserLog(
			"WARN",
			"FIND_WALLET_BY_ID",
			"FAILED",
			"ID Wallet tidak valid",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_WALLET",
			},
		)
		return nil, errors.New("ID Wallet tidak valid")
	}

	wallet, err := s.WalletsRepository.FindWalletsByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Wallet tidak ditemukan")
		}
		s.publishUserLog(
			"WARN",
			"FIND_WALLET_BY_ID",
			"FAILED",
			"Tidak dapat menemukan Wallet",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_WALLET",
			},
		)

		return nil, fmt.Errorf("gagal mencari Wallet: %w", err)
	}
	s.publishUserLog(
		"INFO",
		"FIND_WALLET_BY_ID",
		"SUCCESS",
		"Berhasil menemukan Wallet",
		&userID,
		map[string]any{
			"reason": "SUCCESS_FIND_WALLET",
		},
	)

	return wallet, nil
}

// FindWalletsByGroupID digunakan untuk mengambil semua Wallet dalam Group.
func (s *WalletService) FindWalletsByGroupID(
	groupID string,
	userID int64,
) ([]model.Wallet, error) {

	if strings.TrimSpace(groupID) == "" {
		s.publishUserLog(
			"WARN",
			"FIND_WALLET_BY_ID",
			"FAILED",
			"ID Wallet tidak valid",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_ID_GROUP_WALLET",
			},
		)
		return nil, errors.New("Group ID tidak boleh kosong")
	}

	wallets, err := s.WalletsRepository.FindwalletsByGroupId(groupID)
	if err != nil {
		s.publishUserLog(
			"WARN",
			"FIND_WALLET_BY_ID",
			"FAILED",
			"Gagal Menemuka wallets berdasarkan Group Pengeluaran.",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_WALLET_FROM_GROUP",
			},
		)
		return nil, fmt.Errorf("gagal mengambil Wallet berdasarkan Group: %w", err)
	}
	s.publishUserLog(
		"INFO",
		"FIND_WALLET_BY_GROUPID",
		"SUCCESS",
		"Berhasil menemukan Wallet",
		&userID,
		map[string]any{
			"reason": "SUCCESS_FIND_WALLET",
		},
	)

	return wallets, nil
}

// FindWalletsByUserIDAndGroupID digunakan untuk mencari Wallet
// berdasarkan User ID dan Group ID.
func (s *WalletService) FindWalletsByUserIDAndGroupID(
	userID int64,
	groupID string,
) ([]model.Wallet, error) {

	if userID <= 0 {
		s.publishUserLog(
			"WARN",
			"FIND_WALLET_BY_ID_AND_GROUPID",
			"FAILED",
			"ID Wallet tidak valid",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_ID_GROUP_WALLET",
			},
		)
		return nil, errors.New("User ID tidak valid")
	}

	if strings.TrimSpace(groupID) == "" {
		s.publishUserLog(
			"WARN",
			"FIND_WALLET_BY_ID_AND_GROUPID",
			"FAILED",
			"ID Group tidak valid",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_ID_GROUP_WALLET",
			},
		)
		return nil, errors.New("Group ID tidak boleh kosong")
	}

	wallets, err := s.WalletsRepository.FindByuserIdAndGrupId(
		userID,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal mencari Wallet: %w", err)
	}
	s.publishUserLog(
		"INFO",
		"FIND_WALLET_BY_ID_AND_GROUPID",
		"SUCCESS",
		"Berhasil menemukan Wallet",
		&userID,
		map[string]any{
			"reason": "SUCCESS_FIND_WALLET_BY_ID_AND_GROUPID",
		},
	)

	return wallets, nil
}

// FindWalletByName digunakan untuk mencari Wallet berdasarkan nama.
func (s *WalletService) FindWalletByName(
	walletName string,
	userID int64,
) (*model.Wallet, error) {

	if strings.TrimSpace(walletName) == "" {
		s.publishUserLog(
			"WARN",
			"FIND_WALLET_BY_NAME",
			"FAILED",
			"Name Wallet tidak valid",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_ID_GROUP_WALLET",
			},
		)
		return nil, errors.New("nama Wallet tidak boleh kosong")
	}

	wallet, err := s.WalletsRepository.FindByWalletsName(walletName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Wallet tidak ditemukan")
		}
		s.publishUserLog(
			"WARN",
			"FIND_WALLET_BY_NAME",
			"FAILED",
			"Gagal mecari wallet berdasarkan nama",
			&userID,
			map[string]any{
				"reason": "FAILED_FIND_WALLET_BY_NAME",
			},
		)

		return nil, fmt.Errorf("gagal mencari Wallet: %w", err)
	}
	s.publishUserLog(
		"INFO",
		"FIND_WALLET_BY_NAME",
		"SUCCESS",
		"Berhasil mecari wallet berdasarkan nama",
		&userID,
		map[string]any{
			"reason": "SUCCESS_FIND_WALLET_BY_NAME",
		},
	)

	return wallet, nil
}

// UpdateWallet digunakan untuk mengubah data Wallet.
func (s *WalletService) UpdateWallet(
	wallet *model.Wallet,
	userID int64,
) (*model.Wallet, error) {

	if wallet == nil {
		s.publishUserLog(
			"WARN",
			"UPDATED_WALLET",
			"FAILED",
			"Terdapat Data yang kosong ",
			&userID,
			map[string]any{
				"reason": "FAILED_UPDATED_WALLET",
			},
		)
		return nil, errors.New("data Wallet tidak boleh kosong")

	}
	s.publishUserLog(
		"WARN",
		"UPDATED_WALLET",
		"FAILED",
		"Terdapat Data yang kosong ",
		&userID,
		map[string]any{
			"reason": "FAILED_UPDATED_WALLET",
		},
	)
	if wallet.ID <= 0 {
		s.publishUserLog(
			"WARN",
			"UPDATED_WALLET",
			"FAILED",
			"Terdapat Data yang kosong ",
			&userID,
			map[string]any{
				"reason": "FAILED_UPDATED_WALLET",
			},
		)
		return nil, errors.New("ID Wallet tidak valid")
	}

	if strings.TrimSpace(wallet.WalletName) == "" {
		s.publishUserLog(
			"WARN",
			"UPDATED_WALLET",
			"FAILED",
			"Terdapat Data yang kosong ",
			&userID,
			map[string]any{
				"reason": "FAILED_UPDATED_WALLET",
			},
		)
		return nil, errors.New("nama Wallet wajib diisi")
	}

	result, err := s.WalletsRepository.UpdateWallets(wallet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Wallet tidak ditemukan")
		}

		return nil, fmt.Errorf("gagal mengubah Wallet: %w", err)
	}
	s.publishUserLog(
		"INFO",
		"UPDATED_WALLET",
		"SUCCESS",
		"Berhasil mengupdate wallet ",
		&userID,
		map[string]any{
			"reason": "SUCCESS_UPDATED_WALLET",
		},
	)

	return result, nil
}

// DeleteWallet digunakan untuk menghapus Wallet.
func (s *WalletService) DeleteWallet(
	id int64,
	userID int64,
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

func (s *WalletService) publishUserLog(
	level string,
	eventName string,
	status string,
	message string,
	userID *int64,
	metadata map[string]any,
) {
	if s.LoggingPublisher == nil {
		log.Println("Logging Publisher belum tersedia")
		return
	}

	logData := rabbitmq.LoggingMessage{
		ServiceName:  "Group_service",
		UserID:       userID,
		Level:        level,
		EventName:    eventName,
		Status:       status,
		Message:      message,
		ResourceType: stringPointer("user"),
		Metadata:     metadata,
	}

	err := s.LoggingPublisher.PublishLogging(
		context.Background(),
		logData,
	)

	if err != nil {
		log.Println("Gagal mengirim log registrasi:", err)
	}
}

func stringPointer(value string) *string {
	return &value
}
