# Auth Service

Микросервис аутентификации и авторизации для проекта Agromeration.

## Описание

Auth Service предоставляет REST API для:
- Регистрации пользователей
- Аутентификации пользователей
- Обновления JWT токенов
- Выхода пользователей
- Проверки здоровья сервиса

## Технологии

- **Go 1.21+** - основной язык программирования
- **go-chi** - HTTP роутер
- **oapi-codegen** - генерация кода из OpenAPI спецификации
- **JWT** - JSON Web Tokens для аутентификации
- **bcrypt** - хеширование паролей

## Структура проекта

```
services/auth/
├── api/                    # OpenAPI спецификация и сгенерированный код
│   ├── auth.yaml          # OpenAPI спецификация
│   ├── cfg.yaml           # Конфигурация oapi-codegen
│   ├── generate.go        # go:generate директива
│   └── auth.gen.go        # Сгенерированный код (не редактировать)
├── cmd/                   # Точка входа приложения
│   └── main.go           # Основной файл приложения
├── internal/              # Внутренняя логика сервиса
│   ├── server.go         # Реализация HTTP обработчиков
│   ├── storage.go        # Хранилище данных (в памяти)
│   ├── jwt.go            # JWT сервис
│   ├── password.go       # Сервис для работы с паролями
│   └── errors.go         # Определения ошибок
├── go.mod                # Зависимости Go
├── go.sum                # Хеши зависимостей
├── tools.go              # Зависимости для инструментов
└── README.md             # Этот файл
```

## Установка и запуск

### Предварительные требования

- Go 1.21 или выше
- Git

### Установка зависимостей

```bash
cd services/auth
go mod tidy
```

### Генерация кода

```bash
cd api
go generate
```

### Сборка

```bash
go build -o auth-service ./cmd
```

### Запуск

```bash
./auth-service
```

По умолчанию сервис запускается на порту 8081.

### Параметры запуска

- `-port` - порт для HTTP сервера (по умолчанию: 8081)
- `-secret` - секретный ключ для JWT (по умолчанию: "your-secret-key-change-in-production")

Пример:
```bash
./auth-service -port 8082 -secret "my-super-secret-key"
```

## API Endpoints

### Health Check

```
GET /health
```

Проверка здоровья сервиса.

**Ответ:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Регистрация пользователя

```
POST /auth/register
```

**Тело запроса:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Ответ (201):**
```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  },
  "message": "User registered successfully"
}
```

### Вход пользователя

```
POST /auth/login
```

**Тело запроса:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Ответ (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### Обновление токена

```
POST /auth/refresh
```

**Тело запроса:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ответ (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

### Выход пользователя

```
POST /auth/logout
```

**Заголовки:**
```
Authorization: Bearer <access_token>
```

**Ответ (200):**
Пустой ответ с кодом 200.

## Тестирование

### С помощью curl

```bash
# Health check
curl http://localhost:8081/health

# Регистрация
curl -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","first_name":"John","last_name":"Doe"}'

# Вход
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Обновление токена
curl -X POST http://localhost:8081/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"your-refresh-token"}'

# Выход
curl -X POST http://localhost:8081/auth/logout \
  -H "Authorization: Bearer your-access-token"
```

### С помощью wget

```bash
# Health check
wget -qO- http://localhost:8081/health

# Регистрация
wget --post-data='{"email":"test@example.com","password":"password123","first_name":"John","last_name":"Doe"}' \
     --header='Content-Type:application/json' \
     -qO- http://localhost:8081/auth/register

# Вход
wget --post-data='{"email":"test@example.com","password":"password123"}' \
     --header='Content-Type:application/json' \
     -qO- http://localhost:8081/auth/login
```

## Особенности реализации

### Хранилище данных

В текущей реализации используется хранилище в памяти. Это означает, что:
- Данные теряются при перезапуске сервиса
- Не подходит для продакшена
- Подходит для разработки и тестирования

### Безопасность

- Пароли хешируются с помощью bcrypt
- JWT токены подписываются с использованием HMAC-SHA256
- Access токены имеют срок действия 1 час
- Refresh токены имеют срок действия 7 дней

### Валидация

Все запросы валидируются против OpenAPI спецификации с помощью middleware.

## Разработка

### Изменение API

1. Отредактируйте файл `api/auth.yaml`
2. Перегенерируйте код:
   ```bash
   cd api
   go generate
   ```
3. Обновите реализацию в `internal/server.go`

### Добавление новых эндпоинтов

1. Добавьте новый путь в `api/auth.yaml`
2. Перегенерируйте код
3. Реализуйте новый метод в `internal/server.go`

## Лицензия

MIT 