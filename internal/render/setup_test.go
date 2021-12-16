package render

import (
	"encoding/gob"
	"github.com/ahmedkhaeld/bookings/internal/config"
	"github.com/ahmedkhaeld/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// get the session inforamtion
	gob.Register(models.Reservation{})
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = false
	session.Cookie.SameSite = http.SameSiteLaxMode

	testApp.Session = session

	// make the app in the render.go points to testApp
	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}
func (tw *myWriter) WriteHeader(w int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
