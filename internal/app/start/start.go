package start

import (
	"net/http"

	"golangproject/internal/app/config"
	"golangproject/internal/deliveries/http_delivery"
	"golangproject/internal/services/auth"
	"golangproject/internal/services/user"
	_ "golangproject/cmd/app/docs"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

type Server struct {
    httpServer *http.Server
}


// Start запускает HTTP сервер
func (s *Server) Start() error {
    return s.httpServer.ListenAndServe()
}

func NewServer(authService *auth.Service, userService *user.Service) *Server {
	cfg := config.LoadConfig()
    router := setupRouter(authService, userService)
    
    return &Server {
        httpServer: &http.Server{
            Addr:    cfg.HTTPPort,
            Handler: router,
        },
	}
}


// setupRouter внутренняя функция для настройки маршрутизатора
func setupRouter(authService *auth.Service, userService *user.Service) *mux.Router {
    router := mux.NewRouter()
    
    authMiddleware := http_delivery.AuthMiddleware(authService, userService)
    loginHandler := http_delivery.LoginHandler(authService, userService)
	registerHandler := http_delivery.RegisterHandler(userService)

    // Public routes
	router.Path("/swagger/{any:.*}").Handler(httpSwagger.WrapHandler)
    router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/register", registerHandler).Methods("POST")  // rest methods
	router.Handle("/profile", authMiddleware(http_delivery.ProfileHandler(userService))).Methods("GET")

    // Protected routes
    protected := router.PathPrefix("/api").Subrouter()
    protected.Use(authMiddleware)
    protected.HandleFunc("/profile", http_delivery.ProfileHandler(userService)).Methods("GET")

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }).Methods("GET")
    
    return router
}
