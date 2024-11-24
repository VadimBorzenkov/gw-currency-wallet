package services

import (
	"errors"
	"fmt"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
)

// GetAllRates - получение всех курсов валют.
func (s *service) GetAllRates() (map[string]float64, error) {
	rates, err := s.currencyClient.GetAllRates()
	if err != nil {
		s.logger.Errorf("Failed to fetch exchange rates: %v", err)
		return nil, err
	}
	return rates, nil
}

// GetRate - получение курса обмена для двух валют.
func (s *service) GetRate(fromCurrency, toCurrency string) (float64, error) {
	rate, err := s.currencyClient.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		s.logger.Errorf("Failed to fetch exchange rate for %s to %s: %v", fromCurrency, toCurrency, err)
		return 0, err
	}
	return rate, nil
}

// ConvertCurrency - конвертация валюты.
func (s *service) ExchangeCurrency(fromCurrency, toCurrency string, amount float64) (float64, error) {
	converted, err := s.currencyClient.ConvertCurrency(fromCurrency, toCurrency, amount)
	if err != nil {
		s.logger.Errorf("Failed to convert currency from %s to %s: %v", fromCurrency, toCurrency, err)
		return 0, err
	}
	return converted, nil
}

// CreateWallet создаёт новый кошелёк для пользователя.
func (s *service) CreateWallet(wallet *models.Wallet) (int, error) {
	return s.repo.CreateWallet(wallet)
}

// GetBalance возвращает баланс пользователя по валютам.
func (s *service) GetBalance(userID uint64) (map[string]float64, error) {
	// Получаем все записи кошелька пользователя
	wallets, err := s.repo.GetWalletsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Формируем баланс в виде карты
	balances := make(map[string]float64)
	for _, wallet := range wallets {
		balances[wallet.Currency] = wallet.Balance
	}

	return balances, nil
}

// Deposit пополняет кошелёк пользователя в указанной валюте.
func (s *service) Deposit(userID uint64, amount float64, currency string) (map[string]float64, error) {
	// Проверяем корректность суммы
	if amount <= 0 {
		return nil, errors.New("invalid deposit amount")
	}

	// Получаем кошелёк по валюте
	wallet, err := s.repo.GetWalletByUserAndCurrency(userID, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wallet for currency %s: %v", currency, err)
	}
	if wallet == nil {
		return nil, errors.New("wallet for the specified currency does not exist")
	}

	// Обновляем баланс
	newBalance := wallet.Balance + amount
	err = s.repo.UpdateWalletBalance(wallet.ID, newBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Возвращаем текущие балансы пользователя
	return s.GetAllBalances(userID)
}

// Withdraw выводит средства из кошелька пользователя в указанной валюте.
func (s *service) Withdraw(userID uint64, amount float64, currency string) (map[string]float64, error) {
	// Проверяем корректность суммы
	if amount <= 0 {
		return nil, errors.New("invalid withdrawal amount")
	}

	// Получаем кошелёк по валюте
	wallet, err := s.repo.GetWalletByUserAndCurrency(userID, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wallet for currency %s: %v", currency, err)
	}
	if wallet == nil {
		return nil, errors.New("wallet for the specified currency does not exist")
	}

	// Проверяем баланс
	if wallet.Balance < amount {
		return nil, errors.New("insufficient funds")
	}

	// Обновляем баланс
	newBalance := wallet.Balance - amount
	err = s.repo.UpdateWalletBalance(wallet.ID, newBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Возвращаем текущие балансы пользователя
	return s.GetAllBalances(userID)
}

// GetAllBalances возвращает балансы во всех валютах для пользователя.
func (s *service) GetAllBalances(userID uint64) (map[string]float64, error) {
	wallets, err := s.repo.GetWalletsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wallets: %v", err)
	}

	balances := make(map[string]float64)
	for _, wallet := range wallets {
		balances[wallet.Currency] = wallet.Balance
	}

	return balances, nil
}

func (s *service) UpdateUserBalance(userID uint64, fromCurrency, toCurrency string, amount, exchangedAmount float64) (map[string]float64, error) {
	// Получаем кошелёк, откуда списываем средства
	fromWallet, err := s.repo.GetWalletByUserAndCurrency(userID, fromCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet for currency %s: %v", fromCurrency, err)
	}
	if fromWallet == nil || fromWallet.Balance < amount {
		return nil, fmt.Errorf("insufficient funds in %s wallet", fromCurrency)
	}

	// Получаем кошелёк, куда зачисляем средства
	toWallet, err := s.repo.GetWalletByUserAndCurrency(userID, toCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet for currency %s: %v", toCurrency, err)
	}
	if toWallet == nil {
		return nil, fmt.Errorf("wallet for currency %s does not exist", toCurrency)
	}

	// Обновляем баланс кошелька "FromCurrency"
	newFromBalance := fromWallet.Balance - amount
	err = s.repo.UpdateWalletBalance(fromWallet.ID, newFromBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet balance for %s: %v", fromCurrency, err)
	}

	// Обновляем баланс кошелька "ToCurrency"
	newToBalance := toWallet.Balance + exchangedAmount
	err = s.repo.UpdateWalletBalance(toWallet.ID, newToBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet balance for %s: %v", toCurrency, err)
	}

	// Возвращаем новые балансы
	return map[string]float64{
		fromCurrency: newFromBalance,
		toCurrency:   newToBalance,
	}, nil
}
