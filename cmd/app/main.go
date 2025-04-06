package main

import (
	"golangproject/internal/app"
	_ "golangproject/cmd/app/docs"
)

// @title Clean Architecture API
// @version 1.0
// @description This is a sample server for Clean Architecture.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
    app := app.New()
    app.Run()
}