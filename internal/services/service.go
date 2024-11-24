package services

import (
	"context"
	"time"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/grpc"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/repository"
	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/utils"
	"github.com/sirupsen/logrus"
)

// Service определяет методы бизнес-логики.
type Service interface {
	// User methods
	RegisterUser(user *models.User) (int64, error)
	GetUserByID(userID uint64) (*models.User, error)
	AuthenticateUser(username, password string) (*models.User, error)

	// Wallet methods
	CreateWallet(wallet *models.Wallet) (int, error)
	GetBalance(userID uint64) (map[string]float64, error)
	Deposit(userID uint64, amount float64, currency string) (map[string]float64, error)
	Withdraw(userID uint64, amount float64, currency string) (map[string]float64, error)
	UpdateUserBalance(userID uint64, fromCurrency, toCurrency string, amount, exchangedAmount float64) (map[string]float64, error)
	GetAllBalances(userID uint64) (map[string]float64, error)

	// gw-exchanger methods
	GetAllRates() (map[string]float64, error)
	GetRate(fromCurrency, toCurrency string) (float64, error)
	ExchangeCurrency(fromCurrency, toCurrency string, amount float64) (float64, error)

	NewJWT(userId uint64, email, ipAddress, tokenID string) (string, error)
	AccessTTL() time.Duration
	CreateRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error
	GetRefreshTokenModelByID(ctx context.Context, userID uint64, deviceId string) (*models.RefreshToken, error)
	DeleteRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error
}

type service struct {
	repo           repository.Repository
	currencyClient *grpc.CurrencyClient
	tokenManger    utils.Manager
	logger         *logrus.Logger
}

// Новый сервис с зависимостью от клиента валют
func NewService(repo repository.Repository, currencyClient *grpc.CurrencyClient, tokenManger utils.Manager, logger *logrus.Logger) Service {
	return &service{repo: repo, currencyClient: currencyClient, tokenManger: tokenManger, logger: logger}
}
