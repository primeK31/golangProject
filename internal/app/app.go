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
	// Загрузка конфигурации с обработкой ошибок
	cfg := config.LoadConfig()

	// Инициализация подключения к БД
	db, err := connections.ConnectDB(cfg.SQL_DATABASE_URL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	
	// Убрали defer db.Close() - соединение должно жить пока работает приложение

	fmt.Println("Database connected successfully!")

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(db)

	// Инициализация сервисов
	authService := auth.New(cfg.JWTSecret) // Исправили на JWTSecret
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
