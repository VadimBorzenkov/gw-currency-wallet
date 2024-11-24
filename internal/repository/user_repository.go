package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
)

func (r *repo) CreateUser(user *models.User) (int64, error) {
	query := "INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id"
	var userID int64
	err := r.db.QueryRow(query, user.Username, user.Password, user.Email).Scan(&userID)
	return userID, err
}

func (r *repo) GetUserByID(userID uint64) (*models.User, error) {
	query := "SELECT id, username, password, email FROM users WHERE id = $1"
	user := &models.User{}
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	return user, err
}

func (r *repo) GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT id, username, password, email FROM users WHERE username = $1"
	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	return user, err
}

// Получение RefreshToken по userID и deviceID
func (a *repo) GetRefreshTokenModelByID(ctx context.Context, userID uint64, deviceID string) (*models.RefreshToken, error) {
	query := "SELECT user_id, device_id, token, expires_at FROM refresh_tokens WHERE user_id = ? AND device_id = ?"
	row := a.db.QueryRowContext(ctx, query, userID, deviceID)

	var refreshToken models.RefreshToken
	if err := row.Scan(&refreshToken.UserID, &refreshToken.DeviceID, &refreshToken.Token, &refreshToken.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Не найдено
		}
		a.logger.Error("Error fetching refresh token:", err)
		return nil, err
	}
	return &refreshToken, nil
}

// Добавление нового RefreshToken
func (a *repo) SetRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := "INSERT INTO refresh_tokens (user_id, device_id, token, expires_at) VALUES (?, ?, ?, ?)"
	_, err := a.db.ExecContext(ctx, query, refreshToken.UserID, refreshToken.DeviceID, refreshToken.Token, refreshToken.ExpiresAt)
	if err != nil {
		a.logger.Error("Error inserting refresh token:", err)
		return err
	}
	return nil
}

// Удаление RefreshToken
func (a *repo) DeleteRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := "DELETE FROM refresh_tokens WHERE user_id = ? AND device_id = ?"
	_, err := a.db.ExecContext(ctx, query, refreshToken.UserID, refreshToken.DeviceID)
	if err != nil {
		a.logger.Error("Error deleting refresh token:", err)
		return err
	}
	return nil
}
