package main

import "Notes/internal/app"

// @title Notes API
// @version 1.0
// @description API for Notes application
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	app.Run()
}
