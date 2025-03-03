package server

import (
	"database/sql"
	"errors"
	"git-server/internal/storage"
	"git-server/internal/storage/sqlstorage"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func SetupLogger(env string) (*slog.Logger, error) {
	var logger *slog.Logger
	switch env {
	case EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvDev:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return nil, errors.New("invalid env variable")
	}
	return logger, nil
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

type Server struct {
	config  *Config
	logger  *slog.Logger
	router  *mux.Router
	storage storage.Storage
}

func NewServerWithDB(config *Config) (*Server, error) {
	log, err := SetupLogger(config.Env)
	if err != nil {
		return nil, err
	}
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	s := Server{
		config:  config,
		logger:  log,
		router:  mux.NewRouter(),
		storage: sqlstorage.New(db),
	}

	s.configureRouter()

	return &s, nil
}

// TODO :
// func NewInMemoryServer(config *Config) (*Server, error) {
// 	log, err := SetupLogger(config.Env)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Server{
// 		config:  config,
// 		logger:  log,
// 		router:  mux.NewRouter(),
// 		storage: teststorage.New(),
// 	}, nil
// }

func (server *Server) RunServer() error {
	server.logger.Debug("Server is started")

	return http.ListenAndServe(server.config.Port, server.router)
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func (server *Server) configureRouter() {
	server.router.HandleFunc("/users", server.handleUsersCreate()).Methods("POST")
}

func (server *Server) handleUsersCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// TODO : func Start
