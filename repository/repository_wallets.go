package repository

import (
	"database/sql"
	"fmt"
	"wallets_service/model"
)

type WalletsRepository struct {
	DB *sql.DB
}

func NewWalletsRepository(db *sql.DB) *WalletsRepository {
	return &WalletsRepository{
		DB: db,
	}
}

func (r *WalletsRepository) CreateWallets(
	wallet *model.Wallet,
) (*model.Wallet, error) {

	query := `
		INSERT INTO wallets(
		user_id, 
		group_id, 
		wallet_name, 
		account_number, 
		initial_balance, 
		current_balance,
		is_shared)
    	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.DB.Exec(
		query,
		wallet.UserID,
		wallet.GroupID,
		wallet.WalletName,
		wallet.AccountNumber,
		wallet.InitialBalance,
		wallet.CurrentBalance,
		wallet.IsShared,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat Wallets: %w", err)
	}

	wallet.ID, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan ID Wallet: %w", err)
	}

	return wallet, nil
}

func (r *WalletsRepository) DeleteWallets(id int64) error {
	query := `
        DELETE FROM wallets
        WHERE id = ?
    `

	_, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *WalletsRepository) FindWalletsByID(
	id int64,
) (*model.Wallet, error) {

	wallet := &model.Wallet{}
	err := r.DB.QueryRow(
		`
		SELECT 
		id, 
		user_id, 
		group_id, 
		wallet_name, 
		account_number, 
		initial_balance, 
		current_balance,
		is_shared
		FROM wallets
		WHERE id = ?
		`,
		id,
	).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.GroupID,
		&wallet.WalletName,
		&wallet.AccountNumber,
		&wallet.InitialBalance,
		&wallet.CurrentBalance,
		&wallet.IsShared,
	)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *WalletsRepository) FindwalletsByGroupId(
	groupId string,
) ([]model.Wallet, error) {

	rows, err := r.DB.Query(
		`
        SELECT id, 
		user_id, 
		group_id, 
		wallet_name, 
		account_number, 
		initial_balance, 
		current_balance,
		is_shared
		FROM wallets
        WHERE group_id = ?
        ORDER BY id ASC
        `,
		groupId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallets := []model.Wallet{}

	for rows.Next() {
		var wallet model.Wallet

		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.GroupID,
			&wallet.WalletName,
			&wallet.AccountNumber,
			&wallet.InitialBalance,
			&wallet.CurrentBalance,
			&wallet.IsShared,
		)
		if err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (r *WalletsRepository) FindByWalletsName(
	walletName string,
) (*model.Wallet, error) {

	wallet := &model.Wallet{}
	err := r.DB.QueryRow(
		`
		SELECT id, 
		user_id, 
		group_id, 
		wallet_name, 
		account_number, 
		initial_balance, 
		current_balance,
		is_shared
		FROM wallets
        WHERE wallet_name= ?
		`,
		walletName,
	).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.GroupID,
		&wallet.WalletName,
		&wallet.AccountNumber,
		&wallet.InitialBalance,
		&wallet.CurrentBalance,
		&wallet.IsShared,
	)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (r *WalletsRepository) FindByuserIdAndGrupId(
	userId int64,
	groupId string,
) ([]model.Wallet, error) {

	rows, err := r.DB.Query(
		`
        SELECT id, 
		user_id, 
		group_id, 
		wallet_name, 
		account_number, 
		initial_balance, 
		current_balance,
		is_shared
		FROM wallets
        WHERE user_id = ? AND group_id = ? 
        ORDER BY id ASC
        `,
		userId, groupId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallets := []model.Wallet{}

	for rows.Next() {
		var wallet model.Wallet

		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.GroupID,
			&wallet.WalletName,
			&wallet.AccountNumber,
			&wallet.InitialBalance,
			&wallet.CurrentBalance,
			&wallet.IsShared,
		)
		if err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (r *WalletsRepository) UpdateWallets(wallet *model.Wallet) (*model.Wallet, error) {
	query := `
        UPDATE wallets
        SET wallet_name = ?, 
		account_number = ?, 
		initial_balance = ?, 
		current_balance = ?,
		is_shared = ?
        WHERE id = ?
    `

	_, err := r.DB.Exec(
		query,
		wallet.WalletName,
		wallet.AccountNumber,
		wallet.InitialBalance,
		wallet.CurrentBalance,
		wallet.IsShared,
		wallet.ID,
	)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}
