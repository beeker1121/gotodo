package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gotodo/api"
	"gotodo/api/config"
	"gotodo/database"

	"github.com/beeker1121/creek"
	"github.com/beeker1121/httprouter"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Parse the API configuration file.
	cfg, err := config.ParseConfigFile("config.json")
	if err != nil {
		panic(err)
	}

	// Get the configuration environment variables.
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBPort = os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPass = os.Getenv("DB_PASS")
	cfg.APIHost = os.Getenv("API_HOST")
	cfg.APIPort = os.Getenv("API_PORT")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	// Create new creek logger with 10 MB max file size.
	logger := log.New(creek.New(cfg.LogFile, 10), "Go Todo API: ", log.Llongfile|log.LstdFlags)
	logger.Printf("Starting Go Todo API server at %s\n", time.Now().UTC().Format(time.RFC3339))

	// Connect to the MySQL database.
	db, err := sql.Open("mysql", cfg.DBUser+":"+cfg.DBPass+"@tcp("+cfg.DBHost+":"+cfg.DBPort+")/"+cfg.DBName+"?parseTime=true")
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// Test database connection.
	if err := db.Ping(); err != nil {
		logger.Fatal(err)
	}

	// Create a new Go Todo database.
	gdb := database.New(db)

	// Create a new API.
	router := httprouter.New()
	api.New(cfg, logger, gdb, router)

	// Create a new HTTP server.
	server := &http.Server{
		Addr:           ":" + cfg.APIPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Running server...")

	// Start the HTTP server.
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
