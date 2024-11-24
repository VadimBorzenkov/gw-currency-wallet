package repository

import (
	"database/sql"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
)

func (r *repo) CreateWallet(wallet *models.Wallet) (int, error) {
	query := "INSERT INTO wallets (user_id, balance, currency) VALUES ($1, $2, $3) RETURNING id"
	var walletID int
	err := r.db.QueryRow(query, wallet.UserID, wallet.Balance, wallet.Currency).Scan(&walletID)
	return walletID, err
}

func (r *repo) GetWalletByID(walletID uint64) (*models.Wallet, error) {
	query := "SELECT id, user_id, balance, currency FROM wallets WHERE id = $1"
	wallet := &models.Wallet{}
	err := r.db.QueryRow(query, walletID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency)
	return wallet, err
}

// GetWalletByUserAndCurrency получает кошелёк пользователя по валюте.
func (r *repo) GetWalletByUserAndCurrency(userID uint64, currency string) (*models.Wallet, error) {
	query := "SELECT id, user_id, balance, currency FROM wallets WHERE user_id = $1 AND currency = $2"
	wallet := &models.Wallet{}
	err := r.db.QueryRow(query, userID, currency).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return wallet, err
}

func (r *repo) UpdateWalletBalance(walletID uint64, balance float64) error {
	query := "UPDATE wallets SET balance = $1 WHERE id = $2"
	_, err := r.db.Exec(query, balance, walletID)
	return err
}

// GetWalletsByUserID получает все кошельки пользователя.
func (r *repo) GetWalletsByUserID(userID uint64) ([]*models.Wallet, error) {
	query := "SELECT id, user_id, balance, currency FROM wallets WHERE user_id = $1"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*models.Wallet
	for rows.Next() {
		wallet := &models.Wallet{}
		err := rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency)
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
