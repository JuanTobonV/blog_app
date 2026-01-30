package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/JuanTobonV/blog_app/internal/service"
	"github.com/JuanTobonV/blog_app/internal/store"
	"github.com/JuanTobonV/blog_app/internal/transport"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/blog_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("✅ Connected to database")

	userStore := store.New(db)

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("Failed to get the JWT_SECRET")
	}

	authService := service.NewAuthService(userStore, jwtSecret)

	authHandler := transport.NewAuthHandler(authService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))

}
