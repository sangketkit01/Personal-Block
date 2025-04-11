package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
	"github.com/sangketkit01/personal-block/internal/config"
	"github.com/sangketkit01/personal-block/internal/driver"
	"github.com/sangketkit01/personal-block/internal/handlers"
	"github.com/sangketkit01/personal-block/internal/helpers"
	"github.com/sangketkit01/personal-block/internal/models"
	"github.com/sangketkit01/personal-block/internal/render"
	"github.com/sangketkit01/personal-block/internal/repository"
)

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

const portNumber = ":8216"

func main() {
	gob.Register(models.User{})

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Load environment variables from .env file at project root
	// Calculate path to root directory from cmd/main.go
	rootDir, err := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), ".."))
	if err != nil {
		errorLog.Println("Error finding project root:", err)
	}
	
	err = godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		// Try another approach if the first one fails
		// This handles both running the binary and using 'go run'
		currentDir, _ := os.Getwd()
		projectRoot := filepath.Dir(currentDir)
		err = godotenv.Load(filepath.Join(projectRoot, ".env"))
		if err != nil {
			errorLog.Println("Warning: Error loading .env file:", err)
			// Continue execution as we'll use default values if needed
		}
	}

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln("cannot create template cache", err)
		return
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "personal_blocks")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := driver.NewDatabase(dsn)
	if err != nil {
		errorLog.Fatalln(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	dbRepo := repository.NewDBRepo(&app, db)
	repository.CreateRepo(dbRepo)

	repo := handlers.NewRepository(&app, dbRepo)
	handlers.NewHandlers(repo)

	helpers.NewHelpers(&app)

	render.NewRenderer(&app)

	server := &http.Server{
		Addr:    portNumber,
		Handler: route(),
	}

	fmt.Println("listening on port:", portNumber)

	if err = server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}