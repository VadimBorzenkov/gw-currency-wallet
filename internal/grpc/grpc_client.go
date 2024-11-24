package grpc

import (
	"context"
	"fmt"
	"os"

	exchange_grpc "github.com/VadimBorzenkov/proto-exchange/exchange"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type CurrencyClient struct {
	client exchange_grpc.ExchangeServiceClient
}

// NewCurrencyClient создает новый gRPC клиент с параметрами из конфигурации
func NewCurrencyClient() (*CurrencyClient, error) {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file")
	}

	// Читаем значения переменных окружения для подключения
	grpcAddress := fmt.Sprintf("%s:%s", os.Getenv("GRPC_EXCHANGE_HOST"), os.Getenv("GRPC_EXCHANGE_PORT"))
	if grpcAddress == ":" {
		return nil, fmt.Errorf("Invalid gRPC address")
	}

	// Подключаемся к gRPC серверу
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure()) // Без SSL для простоты
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to gRPC server: %v", err)
	}

	client := exchange_grpc.NewExchangeServiceClient(conn)
	return &CurrencyClient{client: client}, nil
}

// GetExchangeRate возвращает курс обмена между двумя валютами.
func (c *CurrencyClient) GetExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	req := &exchange_grpc.CurrencyRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
	}
	res, err := c.client.GetExchangeRateForCurrency(context.Background(), req)
	if err != nil {
		return 0, err
	}

	return float64(res.GetRate()), nil
}

// GetAllRates возвращает все курсы валют в виде мапы.
func (c *CurrencyClient) GetAllRates() (map[string]float64, error) {
	req := &exchange_grpc.Empty{}
	res, err := c.client.GetExchangeRates(context.Background(), req)
	if err != nil {
		return nil, err
	}

	// Конвертируем курсы в мапу.
	rates := make(map[string]float64)
	for key, value := range res.Rates {
		rates[key] = float64(value)
	}

	return rates, nil
}

// ConvertCurrency конвертирует сумму из одной валюты в другую.
func (c *CurrencyClient) ConvertCurrency(fromCurrency, toCurrency string, amount float64) (float64, error) {
	req := &exchange_grpc.ExchangeRequest{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Amount:       float64(amount),
	}
	res, err := c.client.ConvertCurrency(context.Background(), req)
	if err != nil {
		return 0, err
	}

	return res.ConvertedAmount, nil
}
