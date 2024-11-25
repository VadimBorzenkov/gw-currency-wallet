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
func (r *repo) GetRefreshTokenModelByID(ctx context.Context, userID uint64, deviceID string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, device_id, token, created_at, expires_at 
		FROM refresh_tokens 
		WHERE user_id = $1 AND device_id = $2`
	row := r.db.QueryRowContext(ctx, query, userID, deviceID)

	var refreshToken models.RefreshToken
	if err := row.Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.DeviceID,
		&refreshToken.Token,
		&refreshToken.CreatedAt,
		&refreshToken.ExpiresAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Error fetching refresh token:", err)
		return nil, err
	}
	return &refreshToken, nil
}

// Добавление нового RefreshToken
func (r *repo) SetRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, device_id, token, created_at, expires_at) 
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		refreshToken.UserID,
		refreshToken.DeviceID,
		refreshToken.Token,
		refreshToken.CreatedAt,
		refreshToken.ExpiresAt,
	)
	if err != nil {
		r.logger.Error("Error inserting refresh token:", err)
		return err
	}
	return nil
}

// Удаление RefreshToken
func (r *repo) DeleteRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := `
		DELETE FROM refresh_tokens 
		WHERE user_id = $1 AND device_id = $2`
	_, err := r.db.ExecContext(ctx, query, refreshToken.UserID, refreshToken.DeviceID)
	if err != nil {
		r.logger.Error("Error deleting refresh token:", err)
		return err
	}
	return nil
}
