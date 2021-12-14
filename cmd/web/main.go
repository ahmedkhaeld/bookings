package main

import (
	"fmt"
	"github.com/ahmedkhaeld/bookings/pkg/config"
	"github.com/ahmedkhaeld/bookings/pkg/handlers"
	"github.com/ahmedkhaeld/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

// app entry point to access AppConfig
var app config.AppConfig

// declare the session var
var session *scs.SessionManager

func main() {
	// change this to true when in production environment
	app.InProduction = false

	// initialize the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode

	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	// assign the cache to the app field
	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting the application %v", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
