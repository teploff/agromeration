package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"agromeration/services/auth/api"
)

// AuthServer представляет сервер аутентификации
type AuthServer struct {
	userStorage     *UserStorage
	tokenStorage    *TokenStorage
	jwtService      *JWTService
	passwordService *PasswordService
}

// NewAuthServer создает новый сервер аутентификации
func NewAuthServer(userStorage *UserStorage, secretKey string) *AuthServer {
	return &AuthServer{
		userStorage:     userStorage,
		tokenStorage:    NewTokenStorage(),
		jwtService:      NewJWTService(secretKey),
		passwordService: NewPasswordService(),
	}
}

// GetHealth обрабатывает запрос на проверку здоровья сервиса
func (s *AuthServer) GetHealth(w http.ResponseWriter, r *http.Request) {
	response := api.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Register обрабатывает запрос на регистрацию пользователя
func (s *AuthServer) Register(w http.ResponseWriter, r *http.Request) {
	var req api.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Проверяем, что пользователь не существует
	if _, err := s.userStorage.GetUserByEmail(context.Background(), string(req.Email)); err == nil {
		s.writeError(w, http.StatusConflict, "User already exists", "USER_ALREADY_EXISTS")
		return
	}

	// Создаем пользователя (в реальном приложении здесь нужно хешировать пароль)
	user, err := s.userStorage.CreateUser(context.Background(), string(req.Email), req.Password, req.FirstName, req.LastName)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to create user", "INTERNAL_ERROR")
		return
	}

	response := api.RegisterResponse{
		User:    *user,
		Message: "User registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login обрабатывает запрос на вход пользователя
func (s *AuthServer) Login(w http.ResponseWriter, r *http.Request) {
	var req api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Получаем пользователя по email
	user, err := s.userStorage.GetUserByEmail(context.Background(), string(req.Email))
	if err != nil {
		s.writeError(w, http.StatusUnauthorized, "Invalid credentials", "INVALID_CREDENTIALS")
		return
	}

	// В реальном приложении здесь нужно проверять хеш пароля
	// Для демонстрации просто проверяем, что пароль не пустой
	if req.Password == "" {
		s.writeError(w, http.StatusUnauthorized, "Invalid credentials", "INVALID_CREDENTIALS")
		return
	}

	// Генерируем токены
	accessToken, err := s.jwtService.GenerateAccessToken(user.Id.String(), string(user.Email))
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to generate access token", "INTERNAL_ERROR")
		return
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.Id.String(), string(user.Email))
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to generate refresh token", "INTERNAL_ERROR")
		return
	}

	// Сохраняем refresh token
	s.tokenStorage.StoreRefreshToken(refreshToken, user.Id.String())

	response := api.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 час
		TokenType:    "Bearer",
		User:         *user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RefreshToken обрабатывает запрос на обновление токена
func (s *AuthServer) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req api.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Получаем ID пользователя по refresh token
	userID, err := s.tokenStorage.GetUserIDByRefreshToken(req.RefreshToken)
	if err != nil {
		s.writeError(w, http.StatusUnauthorized, "Invalid refresh token", "INVALID_REFRESH_TOKEN")
		return
	}

	// Получаем пользователя
	user, err := s.userStorage.GetUserByID(context.Background(), userID)
	if err != nil {
		s.writeError(w, http.StatusUnauthorized, "User not found", "USER_NOT_FOUND")
		return
	}

	// Генерируем новый access token
	accessToken, err := s.jwtService.GenerateAccessToken(user.Id.String(), string(user.Email))
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to generate access token", "INTERNAL_ERROR")
		return
	}

	response := api.RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   3600, // 1 час
		TokenType:   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Logout обрабатывает запрос на выход пользователя
func (s *AuthServer) Logout(w http.ResponseWriter, r *http.Request) {
	// Извлекаем токен из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		s.writeError(w, http.StatusUnauthorized, "Authorization header required", "MISSING_AUTHORIZATION")
		return
	}

	// Проверяем формат заголовка
	if !strings.HasPrefix(authHeader, "Bearer ") {
		s.writeError(w, http.StatusUnauthorized, "Invalid authorization header format", "INVALID_AUTHORIZATION_FORMAT")
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Валидируем токен
	_, err := s.jwtService.ValidateToken(token)
	if err != nil {
		s.writeError(w, http.StatusUnauthorized, "Invalid token", "INVALID_TOKEN")
		return
	}

	// В реальном приложении здесь можно добавить токен в черный список
	// Для демонстрации просто возвращаем успешный ответ

	w.WriteHeader(http.StatusOK)
}

// writeError записывает ошибку в ответ
func (s *AuthServer) writeError(w http.ResponseWriter, statusCode int, message, code string) {
	errorResponse := api.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
