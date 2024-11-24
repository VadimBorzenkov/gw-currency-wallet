package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
)

func (s *service) NewJWT(userId uint64, email, ipAddress, tokenID string) (string, error) {
	return s.tokenManger.NewJWT(userId, email, ipAddress, tokenID)
}

func (s *service) AccessTTL() time.Duration {
	return s.tokenManger.GetAccessTTL()
}

func (s *service) CreateRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error {
	return s.repo.SetRefreshTokenModel(ctx, refreshToken)

}

func (s *service) GetRefreshTokenModelByID(ctx context.Context, userID uint64, deviceId string) (*models.RefreshToken, error) {
	return s.repo.GetRefreshTokenModelByID(ctx, userID, deviceId)
}

func (s *service) DeleteRefreshTokenModel(ctx context.Context, refreshToken *models.RefreshToken) error {
	return s.repo.DeleteRefreshTokenModel(ctx, refreshToken)
}

// RegisterUser регистрирует нового пользователя.
func (s *service) RegisterUser(user *models.User) (int64, error) {
	// Проверка на существование пользователя с таким же именем.
	existingUser, err := s.repo.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return 0, errors.New("username already exists")
	}

	// Хэширование пароля.
	hashedPassword, err := s.tokenManger.HashPassword(string(user.Password))
	if err != nil {
		log.Println("Failed to hash password:", err)
		return 0, errors.New("internal server error")
	}
	user.Password = hashedPassword

	// Создание пользователя в базе данных.
	return s.repo.CreateUser(user)
}

// GetUserByID возвращает пользователя по его ID.
func (s *service) GetUserByID(userID uint64) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

// AuthenticateUser аутентифицирует пользователя.
func (s *service) AuthenticateUser(username, password string) (*models.User, error) {
	// Получение пользователя из базы данных.
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Валидация пароля.
	err = s.tokenManger.ValidatePassword(password, user.Password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
