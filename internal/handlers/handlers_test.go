package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// postData holds the posted info its key and value
type postData struct {
	key   string
	value string
}

// theTests slice of structs, each struct has those fields to be tested
var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	// table test for the theTests
	for _, individualTest := range theTests {
		if individualTest.method == "GET" {
			// make a request as a client as web browser accessing a web page
			resp, err := testServer.Client().Get(testServer.URL + individualTest.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != individualTest.expectedStatusCode {
				t.Errorf("for %s, expected %d but go %d", individualTest.name, individualTest.expectedStatusCode, resp.StatusCode)
			}
		} else {

		}
	}

}
