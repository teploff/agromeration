package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"agromeration/services/auth/api"
	"agromeration/services/auth/internal"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	port := flag.String("port", "8081", "Port for auth HTTP server")
	secretKey := flag.String("secret", "your-secret-key-change-in-production", "JWT secret key")
	pgConn := flag.String("pg", os.Getenv("AUTH_PG_DSN"), "PostgreSQL DSN")
	flag.Parse()

	// Подключение к PostgreSQL
	pool, err := pgxpool.New(context.Background(), *pgConn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()

	// Загружаем OpenAPI спецификацию
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Очищаем массив серверов в swagger спецификации, чтобы пропустить валидацию
	// соответствия имен серверов. Мы не знаем, как это будет запущено.
	swagger.Servers = nil

	// Создаем UserStorage с поддержкой PostgreSQL
	userStorage := internal.NewUserStorage(pool)
	// Создаем экземпляр нашего обработчика, который удовлетворяет сгенерированному интерфейсу
	authServer := internal.NewAuthServer(userStorage, *secretKey)

	// Настраиваем базовый chi роутер
	r := chi.NewRouter()

	// Используем middleware валидации для проверки всех запросов против OpenAPI схемы
	r.Use(middleware.OapiRequestValidator(swagger))

	// Добавляем CORS middleware для разработки
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Регистрируем наш authServer как обработчик для интерфейса
	api.HandlerFromMux(authServer, r)

	// Создаем HTTP сервер
	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", *port),
	}

	// Грейсфул-шатдаун
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	// Логируем информацию о запуске
	log.Printf("Auth service starting on port %s", *port)
	log.Printf("Health check available at: http://localhost:%s/health", *port)
	log.Printf("API documentation available at: http://localhost:%s/docs", *port)

	// Запускаем HTTP сервер до конца света
	log.Fatal(s.ListenAndServe())
}
