package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/ahmedkhaeld/bookings/internal/config"
	"github.com/ahmedkhaeld/bookings/internal/models"
	"github.com/ahmedkhaeld/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var functions = template.FuncMap{}
var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "../../templates"

func TestMain(m *testing.M) {
	// #copy of Run() in main
	// 1. register the models' reservation into the sessions
	gob.Register(models.Reservation{})

	// 2. set to false
	app.InProduction = false

	infoLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// 3. create the session then store it in a var
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode

	app.Session = session

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)

	listenForMail()

	// 4. create the template cache
	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")

	}

	// assign the cache to the app field
	app.TemplateCache = templateCache
	// 5. set UseCache to true to render from the cache instead render from disk
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

// getRoutes return routes http handlers
func getRoutes() http.Handler {

	// #copy of routes() in routes.go

	mux := chi.NewRouter()

	// Gracefully absorb panics and prints the stack trace
	mux.Use(middleware.Recoverer)

	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("../../static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

// NoSurf adds CSRF protection against all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// #copy of middleware file
// SessionLoad loads and saves the session on every request for a middleware
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	//filepath.Glob get the location of template pages.
	pagesPath, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return cache, err
	}

	// for loop extract the page name for the pages' path.
	for _, page := range pagesPath {
		pageName := filepath.Base(page)

		templateSet, err := template.New(pageName).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}
		// check template matches any layouts
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
			if err != nil {
				return cache, err
			}
		}
		cache[pageName] = templateSet
	}
	return cache, nil
}
