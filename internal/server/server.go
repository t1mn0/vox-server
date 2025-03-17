package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"vox-server/internal/models"
	"vox-server/internal/storage"
	"vox-server/internal/storage/postgres_storage"
	"vox-server/internal/storage/test_storage"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type contextKey int16

const (
	userContextKey contextKey = iota
	requestIDContextKey
)

type Server struct {
	config  *Config
	logger  *slog.Logger
	router  *mux.Router
	storage storage.Storage
}

func initDB(database_url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", database_url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewServerWithDB(config *Config, useTestDB bool) (*Server, error) {
	log, err := SetupLogger(config.Env)
	if err != nil {
		return nil, err
	}

	var db *sql.DB
	if useTestDB {
		db, err = initDB(config.TestDatabaseURL)
	} else {
		db, err = initDB(config.DatabaseURL)
	}

	if err != nil {
		return nil, err
	}

	s := Server{
		config:  config,
		logger:  log,
		router:  mux.NewRouter(),
		storage: postgres_storage.NewDBStorage(db),
	}

	s.configureRouter()

	return &s, nil
}

func NewInMemoryServer(config *Config) (*Server, error) {
	log, err := SetupLogger(config.Env)
	if err != nil {
		return nil, err
	}

	return &Server{
		config:  config,
		logger:  log,
		router:  mux.NewRouter(),
		storage: test_storage.NewInMemoryStorage(),
	}, nil
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func (server *Server) configureRouter() {
	server.router.Use(server.setRequestID)
	server.router.Use(server.logRequest)
	server.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"}))) // any domain can make requests to your server
	server.router.HandleFunc("/users", server.handleUsersCreate()).Methods("POST")
	server.router.HandleFunc("/sessions", server.handleSessionsCreate()).Methods("POST")

	private := server.router.PathPrefix("/private").Subrouter()
	private.Use(server.authentificateUser)
	private.HandleFunc("/whoami", server.handleWhoAmI()).Methods("GET")
}

func (server *Server) RunServer() error {
	server.logger.Debug("Server is started")
	return http.ListenAndServe(server.config.Port, server.router)
}

func StartServer(useTestDB bool) (*Server, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	if err = ConfigurationDBs(cfg); err != nil {
		return nil, err
	}

	s, err := NewServerWithDB(cfg, useTestDB)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (server *Server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDContextKey, id)))
	})
}

func (server *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		end := time.Now()

		log.Printf(
			"\nREQUEST:\n    req_id=%s\n    method=%s\n    status_code=%d\n    status_msg=%s\n    path=%s\n    remote_addr=%s\n    start=%s\n    end=%s\n    took=%v\n",
			r.Context().Value(requestIDContextKey),
			r.Method,
			rw.code,
			http.StatusText(rw.code),
			r.URL.Path,
			r.RemoteAddr,
			start,
			end,
			end.Sub(start),
		)
	})
}

func (server *Server) authentificateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// if there is no authorization header, skip ahead (the route may be publicly accessible)
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			server.error(w, r, http.StatusUnauthorized, fmt.Errorf("invalid authorization header format"))
			return
		}

		tokenString := parts[1]

		claims, err := ValidateToken(tokenString)
		if err != nil {
			server.error(w, r, http.StatusUnauthorized, fmt.Errorf("invalid token: %w", err))
			return
		}

		var u *models.User
		if strings.Contains(claims.LoginOrEmail, "@") {
			u, err = server.storage.Users().FindByEmail(claims.LoginOrEmail)
		} else {
			u, err = server.storage.Users().FindByLogin(claims.LoginOrEmail)
		}

		if err != nil {
			server.error(w, r, http.StatusUnauthorized, fmt.Errorf("invalid token: %w", err))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, u)))
	})
}

func (server *Server) handleWhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userContextKey).(*models.User)
		if !ok || user == nil {
			server.error(w, r, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		user.Sanitize()
		server.respond(w, r, http.StatusOK, user)
	}
}

func (server *Server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &models.User{
			Login:    req.Login,
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		if err := server.storage.Users().Create(u); err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
		}

		u.Sanitize()

		accessToken, refreshToken, err := GenerateToken(req.Login)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, fmt.Errorf("failed to generate token: %w", err))
			return
		}

		response := map[string]any{
			"user":          u,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}

		server.respond(w, r, http.StatusCreated, response)
	}
}

// check authen
func (server *Server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		LoginOrEmail string `json:"login_or_email"`
		Password     string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.LoginOrEmail == "" || req.Password == "" {
			server.error(w, r, http.StatusBadRequest, errors.New("login/email and password are required"))
			return
		}

		var u *models.User
		var err error

		if strings.Contains(req.LoginOrEmail, "@") {
			u, err = server.storage.Users().FindByEmail(req.LoginOrEmail)
		} else {
			u, err = server.storage.Users().FindByLogin(req.LoginOrEmail)
		}

		if err != nil {
			server.error(w, r, http.StatusUnauthorized, errors.New("incorrect login/email or password"))
			return
		}

		if !u.ComparePassword(req.Password) {
			server.error(w, r, http.StatusUnauthorized, errors.New("incorrect login/email or password"))
			return
		}

		accessToken, refreshToken, err := GenerateToken(req.LoginOrEmail)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, fmt.Errorf("failed to generate token: %w", err))
			return
		}

		response := map[string]any{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}

		server.respond(w, r, http.StatusOK, response)
	}
}

// Render error
func (server *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	server.respond(w, r, code, map[string]string{"error": err.Error()})
}

// Render all types feedback
func (server *Server) respond(w http.ResponseWriter, _ *http.Request, code int, data any) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
