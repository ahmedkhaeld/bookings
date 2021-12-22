package render

import (
	"github.com/ahmedkhaeld/bookings/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSessionData()
	if err != nil {
		t.Error(err)
	}

	/// populate the session context
	session.Put(r.Context(), "flash", "oops!")

	result := AddDefaultData(&td, r)
	if result.Flash != "oops!" {
		t.Error("flash value not found in session")
	}
}

func TestTemplate(t *testing.T) {

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc
	// to get httpRequest
	r, err := getSessionData()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&ww, r, "not-exist.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}

// getSession builds a http request that has session data(context)
func getSessionData() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	// get the context from the request
	ctx := r.Context()
	// put session data in ctx
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	// put the context back into the request
	r = r.WithContext(ctx)

	return r, nil
}
