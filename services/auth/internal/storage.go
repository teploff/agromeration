package internal

import (
	"sync"

	"agromeration/services/auth/api"
	"agromeration/services/auth/internal/db"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oapi-codegen/runtime/types"
)

// UserStorage теперь работает с PostgreSQL через sqlc
// users map/email больше не используется
// mu не нужен

type UserStorage struct {
	q *db.Queries
}

// TokenStorage представляет хранилище токенов в памяти
type TokenStorage struct {
	refreshTokens map[string]string // refresh_token -> user_id
	mu            sync.RWMutex
}

// NewUserStorage создает новое хранилище пользователей с подключением к БД
func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		q: db.New(pool),
	}
}

// NewTokenStorage создает новое хранилище токенов
func NewTokenStorage() *TokenStorage {
	return &TokenStorage{
		refreshTokens: make(map[string]string),
	}
}

// CreateUser создает нового пользователя в БД
func (s *UserStorage) CreateUser(ctx context.Context, email, password, firstName, lastName string) (*api.User, error) {
	// username = email, firstName/lastName не сохраняются (можно расширить схему позже)
	user, err := s.q.CreateUser(ctx, db.CreateUserParams{
		Username:     email,
		PasswordHash: password,
		Email:        email,
	})
	if err != nil {
		return nil, ErrUserAlreadyExists // можно обработать ошибку уникальности
	}
	return dbUserToAPIUser(&user), nil
}

// GetUserByEmail получает пользователя по email из БД
func (s *UserStorage) GetUserByEmail(ctx context.Context, email string) (*api.User, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return dbUserToAPIUser(&user), nil
}

// GetUserByID получает пользователя по ID (int32 -> string)
func (s *UserStorage) GetUserByID(ctx context.Context, id string) (*api.User, error) {
	// В текущей схеме id int32, а не UUID. Можно доработать схему позже.
	return nil, ErrUserNotFound // не реализовано
}

// StoreRefreshToken сохраняет refresh token
func (s *TokenStorage) StoreRefreshToken(refreshToken, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refreshTokens[refreshToken] = userID
}

// GetUserIDByRefreshToken получает ID пользователя по refresh token
func (s *TokenStorage) GetUserIDByRefreshToken(refreshToken string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.refreshTokens[refreshToken]
	if !exists {
		return "", ErrInvalidRefreshToken
	}

	return userID, nil
}

// RemoveRefreshToken удаляет refresh token
func (s *TokenStorage) RemoveRefreshToken(refreshToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.refreshTokens, refreshToken)
}

// dbUserToAPIUser преобразует db.User в api.User
func dbUserToAPIUser(u *db.User) *api.User {
	return &api.User{
		Id:        types.UUID(uuid.New()), // Временно генерируем новый UUID, т.к. в БД int32
		Email:     types.Email(u.Email),
		FirstName: "", // не хранится
		LastName:  "", // не хранится
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.CreatedAt, // updated_at не хранится
	}
}
