package app

import (
	"fmt"
	"log"

	"golangproject/internal/app/config"
	"golangproject/internal/app/connections"
	"golangproject/internal/app/start"
	"golangproject/internal/repositories"
	"golangproject/internal/services/auth"
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

	// repo init
	userRepo := repositories.NewUserRepository(db)

	// service init
	authService := auth.New(cfg.JWTSecret)
	userService := user.New(userRepo)
    
    return &App{
        Config: cfg,
        Server: start.NewServer(authService, userService),
    }
}

func (a *App) Run() {
	defer connections.CloseDB()
    if err := a.Server.Start(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
