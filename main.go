package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Unique-GIT/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	// Get DB Queries
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return
	}
	dbQueries := database.New(db)

	// Get Platform
	platform := os.Getenv("PLATFORM")

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))

	config := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	router := http.NewServeMux()
	router.Handle("/app/",
		middlewareLog(
			config.middlewareIncrementMetrics(
				fileHandler,
			),
		),
	)
	router.Handle("POST /api/users",
		middlewareLog(
			config.middlewareIncrementMetrics(
				http.HandlerFunc(config.handlerUser),
			),
		),
	)
	router.Handle("POST /api/chirps",
		middlewareLog(
			config.middlewareIncrementMetrics(
				http.HandlerFunc(config.validate_chirp),
			),
		),
	)
	router.Handle("GET /api/chirps",
		middlewareLog(
			config.middlewareIncrementMetrics(
				http.HandlerFunc(config.getChirps),
			),
		),
	)
	router.Handle("GET /api/healthz",
		middlewareLog(
			config.middlewareIncrementMetrics(
				http.HandlerFunc(handlerHealth),
			),
		),
	)
	router.Handle("GET /admin/metrics",
		middlewareLog(
			http.HandlerFunc(config.logHits),
		),
	)
	router.Handle("POST /admin/reset",
		middlewareLog(
			http.HandlerFunc(config.handlerDeleteUser),
		),
	)

	server := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	server.ListenAndServe()
}
