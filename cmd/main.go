package main

import (
	"fmt"
	"github.com/sangketkit01/personal-block/internal/driver"
	"github.com/sangketkit01/personal-block/internal/repository"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/sangketkit01/personal-block/internal/config"
	"github.com/sangketkit01/personal-block/internal/handlers"
	"github.com/sangketkit01/personal-block/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

const portNumber = ":8216"

func main() {
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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

	db, err := driver.NewDatabase("host=localhost port=5432 user=postgres password=0627457454New " +
		"dbname=personal_blocks")

	if err != nil {
		errorLog.Fatalln(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	dbRepo := repository.NewDBRepo(&app, db)
	repository.CreateRepo(dbRepo)

	repo := handlers.NewRepository(&app, dbRepo)
	handlers.NewHandlers(repo)

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
