package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Config interface {
	GetAuthJWTPublicKeyPath() string
	GetAuthJWTPrivateKeyPath() string
	GetAccessTokenExpiration() time.Duration
	GetRefreshTokenExpiration() time.Duration
}

type Claims struct {
	UserID    uint64 `json:"id"`
	Email     string `json:"sub"`
	IPAddress string `json:"ip"`
	TokenID   string `json:"tip"`
	jwt.StandardClaims
}

type TokenManager interface {
	NewJWT(userId uint64, email, ipAddress, tokenID string) (string, error)
	ParseJWT(accessToken string) (*Claims, error)
	HashPassword(password string) (string, error)
	ValidatePassword(password, hashedPassword string) error
	GetAccessTTL() time.Duration
}

type Manager struct {
	PublicKey  string
	PrivateKey string
	AccessTTL  time.Duration
}

func NewManager(cfg Config) Manager {
	return Manager{
		PublicKey:  cfg.GetAuthJWTPublicKeyPath(),
		PrivateKey: cfg.GetAuthJWTPrivateKeyPath(),
		AccessTTL:  cfg.GetAccessTokenExpiration(),
	}
}

func (m *Manager) NewJWT(userId uint64, email, ipAddress, tokenID string) (string, error) {
	claims := Claims{
		UserID:    userId,
		Email:     email,
		IPAddress: ipAddress,
		TokenID:   tokenID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.AccessTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKeyData, err := os.ReadFile(m.PrivateKey)
	if err != nil {
		log.Fatalf("could not read private key file: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Println(err)
		log.Println(privateKey)
		return "", fmt.Errorf("could not parse private key: %v", err)
	}
	log.Println(token)
	return token.SignedString(privateKey)
}

func (m *Manager) ParseJWT(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok || token.Header["alg"] != "RS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// fdfdfdf
		log.Println(m.PublicKey)

		publicKeyData, err := os.ReadFile(m.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("could not read private key file: %v", err)
		}
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
		if err != nil {
			return nil, fmt.Errorf("could not parse public key: %v", err)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (m *Manager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), nil
}

func (m *Manager) ValidatePassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (m *Manager) GetAccessTTL() time.Duration {
	return m.AccessTTL
}
