package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ahmedkhaeld/bookings/internal/config"
	"github.com/ahmedkhaeld/bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

//var functions a map of functions that can be used in templates e.g. format a date
// some time we will create our own functions and pass them to the template
var functions = template.FuncMap{
	"humanDate": HumanDate,
}

// HumanDate returns time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// app pointer will have access to the configuration to access TemplateCache or other AppConfig fields
var app *config.AppConfig

// assign template path to variable instead of hard coding
//var pathToTemplates = "../../templates"

// NewRenderer  set app to the AppConfig when it is called to use the TemplateCache
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuth = 1
	}
	return td
}
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var templateCache map[string]*template.Template
	if app.UseCache {
		// get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()

	}

	theTemplate, ok := templateCache[tmpl]
	if !ok {
		return errors.New("can not get template from cache ")
	}

	aBuffer := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	_ = theTemplate.Execute(aBuffer, td)

	_, err := aBuffer.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing to the browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache return a map that has the parsed templates include the layouts
func CreateTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	//filepath.Glob get the location of template pages.
	pagesPath, err := filepath.Glob("../../templates/*.page.tmpl")
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
		matches, err := filepath.Glob("../../templates/*.layout.tmpl")
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("../../templates/*.layout.tmpl")
			if err != nil {
				return cache, err
			}
		}
		cache[pageName] = templateSet
	}
	return cache, nil
}
