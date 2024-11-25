package models

import "time"

type User struct {
	ID           uint64         `json:"id" db:"id"`
	Username     string         `json:"username" db:"username"`
	Password     string         `json:"password" db:"password"`
	Email        string         `json:"email" db:"email"`
	RefreshToken []RefreshToken `json:"refreshToken" db:"refreshToken"`
}

type RefreshToken struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	Token     string    `db:"token"`
	DeviceID  string    `db:"device_id"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

type Wallet struct {
	ID       uint64  `json:"id" db:"id"`
	UserID   uint64  `json:"user_id" db:"user_id"`
	Balance  float64 `json:"balance" db:"balance"`
	Currency string  `json:"currency" db:"currency"`
}

// RegisterRequest представляет тело запроса для регистрации пользователя
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3"` // Логин пользователя
	Password string `json:"password" validate:"required,min=6"` // Пароль
	Email    string `json:"email" validate:"required,email"`    // Электронная почта
}

// RegisterResponse представляет тело ответа при успешной регистрации
type RegisterResponse struct {
	Message string `json:"message"` // Сообщение об успешной регистрации
	UserID  int    `json:"user_id"` // Идентификатор нового пользователя
}

// ErrorResponse ответ ошибки
type ErrorResponse struct {
	Error string `json:"error"`
}

// LoginRequest представляет тело запроса для авторизации
type LoginRequest struct {
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"password123"`
}

// LoginResponse представляет успешный ответ при авторизации
type LoginResponse struct {
	Token string `json:"token" example:"JWT_TOKEN"`
}

// DepositRequest представляет запрос на пополнение баланса.
type DepositRequest struct {
	Amount   float64 `json:"amount" validate:"required,gt=0"`
	Currency string  `json:"currency" validate:"required"`
}

// WithdrawRequest представляет запрос на снятие средств.
type WithdrawRequest struct {
	Amount   float64 `json:"amount" validate:"required,gt=0"`
	Currency string  `json:"currency" validate:"required"`
}

// ExchangeRequest представляет запрос на обмен валют.
type ExchangeRequest struct {
	FromCurrency string  `json:"from_currency" validate:"required"`
	ToCurrency   string  `json:"to_currency" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
}

// BalanceResponse представляет ответ с балансом пользователя.
type BalanceResponse struct {
	Balance map[string]float64 `json:"balance"`
}

// DepositResponse представляет ответ на успешное пополнение баланса.
type DepositResponse struct {
	Message    string  `json:"message"`
	NewBalance float64 `json:"new_balance"`
}

// WithdrawResponse представляет ответ на успешное снятие средств.
type WithdrawResponse struct {
	Message    string  `json:"message"`
	NewBalance float64 `json:"new_balance"`
}

// ExchangeResponse представляет ответ на успешный обмен валют.
type ExchangeResponse struct {
	Message         string             `json:"message"`
	ExchangedAmount float64            `json:"exchanged_amount"`
	NewBalance      map[string]float64 `json:"new_balance"`
}

// RatesResponse представляет ответ с текущими курсами валют.
type RatesResponse struct {
	Rates map[string]float64 `json:"rates"`
	USD   float64            `json:"USD"`
	RUB   float64            `json:"RUB"`
	EUR   float64            `json:"EUR"`
}
