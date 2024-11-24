package repository

import (
	"context"
	"database/sql"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
	"github.com/sirupsen/logrus"
)

// Repository определяет интерфейс для работы с хранилищем данных.
type Repository interface {
	// User methods
	CreateUser(user *models.User) (int64, error)
	GetUserByID(userID uint64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)

	// Wallet methods
	CreateWallet(wallet *models.Wallet) (int, error)
	GetWalletByID(walletID uint64) (*models.Wallet, error)
	UpdateWalletBalance(walletID uint64, balance float64) error
	GetWalletsByUserID(userID uint64) ([]*models.Wallet, error)
	GetWalletByUserAndCurrency(userID uint64, currency string) (*models.Wallet, error)

	// RefreshToken methods
	GetRefreshTokenModelByID(ctx context.Context, userID uint64, deviceID string) (*models.RefreshToken, error)
	SetRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error
	DeleteRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error
}

type repo struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewRepository создаёт новый экземпляр репозитория.
func NewRepository(db *sql.DB, logger *logrus.Logger) Repository {
	return &repo{db: db, logger: logger}
}
