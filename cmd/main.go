package main

import "github.com/VadimBorzenkov/gw-currency-wallet/internal/app"

// @title Currency wallet API
// @version 1.0
// @description API for managing currency wallets and transactions
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	app.Run()

}
