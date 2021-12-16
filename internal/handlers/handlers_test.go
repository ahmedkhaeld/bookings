package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", "POST",
		[]postData{
			{key: "start", value: "2021-12-16"},
			{key: "end", value: "2021-12-17"},
		}, http.StatusOK},

	{"post-search-availability-json", "/search-availability-json", "POST",
		[]postData{
			{key: "start", value: "2021-12-16"},
			{key: "end", value: "2021-12-17"},
		}, http.StatusOK},

	{"post-mr", "/make-reservation", "POST",
		[]postData{
			{key: "first_name", value: "ahmed"},
			{key: "last_name", value: "khaled"},
			{key: "email", value: "hamo@xd.com"},
			{key: "phone", value: "124-3532-5326"},
		}, http.StatusOK},
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
			values := url.Values{}
			for _, x := range individualTest.params {
				values.Add(x.key, x.value)
			}
			resp, err := testServer.Client().PostForm(testServer.URL+individualTest.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != individualTest.expectedStatusCode {
				t.Errorf("for %s, expected %d but go %d", individualTest.name, individualTest.expectedStatusCode, resp.StatusCode)
			}
		}
	}

}
