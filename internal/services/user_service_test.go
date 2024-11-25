package services

import (
	"errors"
	"testing"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/config"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/models"
	"github.com/VadimBorzenkov/gw-currency-wallet/internal/repository/mocks"
	"github.com/VadimBorzenkov/gw-currency-wallet/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg, err := config.LoadConfig()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := logrus.New()

	manager := utils.NewManager(cfg)

	service := NewService(mockRepo, nil, manager, logger)

	user := &models.User{
		Username: "test_user",
		Password: "password",
		Email:    "test@example.com",
	}

	mockRepo.EXPECT().CreateUser(user).Return(int64(1), nil)

	userID, err := service.RegisterUser(user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), userID)
}

func TestAuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаём мок репозитория
	mockRepo := mocks.NewMockRepository(ctrl)

	// Настраиваем зависимые объекты
	cfg, _ := config.LoadConfig()
	logger := logrus.New()
	manager := utils.NewManager(cfg)
	service := NewService(mockRepo, nil, manager, logger)

	// Тестовые данные
	validUser := &models.User{
		ID:       1,
		Username: "testuser",
		Password: "hashedpassword",
	}
	username := "testuser"
	password := "correctpassword"

	// Успешный сценарий
	mockRepo.EXPECT().GetUserByUsername(username).Return(validUser, nil)

	user, err := service.AuthenticateUser(username, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, validUser.Username, user.Username)

	// Сценарий: пользователь не найден
	mockRepo.EXPECT().GetUserByUsername("nonexistent").Return(nil, errors.New("user not found"))

	user, err = service.AuthenticateUser("nonexistent", password)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "user not found")
}
