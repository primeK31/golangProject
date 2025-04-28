package start

import (
	"net/http"

	_ "golangproject/docs/swagger"
	"golangproject/internal/app/config"
	"golangproject/internal/deliveries/http_delivery"
	"golangproject/internal/services/auth"
	"golangproject/internal/services/bet"
	"golangproject/internal/services/session.go"
	"golangproject/internal/services/user"

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

func NewServer(authService *auth.Service, userService *user.Service,
     sessionService *session.Service, betService *bet.SecondServiceClient) *Server {
	cfg := config.LoadConfig()
    router := setupRouter(authService, userService, sessionService, betService)
    
    return &Server {
        httpServer: &http.Server{
            Addr:    cfg.HTTPPort,
            Handler: router,
        },
	}
}

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
func setupRouter(authService *auth.Service, userService *user.Service, 
    sessionService *session.Service, betService *bet.SecondServiceClient) *mux.Router {
    router := mux.NewRouter()
    
    authMiddleware := http_delivery.AuthMiddleware(authService, userService)
    loginHandler := http_delivery.LoginHandler(authService, userService, sessionService)
    logoutHandler := http_delivery.LogoutHandler(sessionService)
	registerHandler := http_delivery.RegisterHandler(userService)
    betHandler := authMiddleware(
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            http_delivery.BetHandler(w, r, betService)
        }),
    )

    // public routes
	router.Path("/swagger/{any:.*}").Handler(httpSwagger.WrapHandler)
    router.HandleFunc("/login", loginHandler).Methods("POST")
    router.Handle("/bet", betHandler).Methods("POST")
    router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/register", registerHandler).Methods("POST")  // rest methods
	router.Handle("/profile", authMiddleware(http_delivery.ProfileHandler(userService))).Methods("GET")

    // protected routes
    protected := router.PathPrefix("/api").Subrouter()
    protected.Use(authMiddleware)
    protected.HandleFunc("/profile", http_delivery.ProfileHandler(userService)).Methods("GET")

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }).Methods("GET")
    
    return router
}
