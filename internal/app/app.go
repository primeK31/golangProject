package app

import (
	"fmt"
	"log"
	//"os"

	"golangproject/internal/app/config"
	"golangproject/internal/app/connections"
	"golangproject/internal/app/start"
	"golangproject/internal/repositories"
	"golangproject/internal/services/auth"
	"golangproject/internal/services/bet"
	"golangproject/internal/services/session.go"
	"golangproject/internal/services/user"
)

type App struct {
	Config *config.Config
	Server *start.Server
}

func New() *App {
	// config load
	cfg := config.LoadConfig()

	// db connect init
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DB_USER, cfg.DB_PASS, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)
	db, err := connections.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("Database connected successfully!")

	apiGatewayURL := cfg.API_GATEWAY_URL
	servicePath := cfg.SERVICE_PATH

	// repo init
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)

	// service init
	authService := auth.New(cfg.JWTSecret)
	userService := user.New(userRepo)
	sessionService := session.New(sessionRepo, cfg.SessionDuration, []byte(cfg.JWTSecret))
	betService := bet.NewSecondServiceClient(apiGatewayURL, servicePath)
    
    return &App{
        Config: cfg,
        Server: start.NewServer(authService, userService, sessionService, betService),
    }
}

func (a *App) Run() {
	defer connections.CloseDB()
    if err := a.Server.Start(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
