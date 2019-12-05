package handlers

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestServeHTTP(t *testing.T) {

	newHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	server := httptest.NewServer(NewCors(newHandler))
	req, _ := http.Get(string(server.URL))

	cases := []struct {
		key           		string
		expectedValue		string
	}{
		{
			"Access-Control-Allow-Origin",
			"*",
		},
		{
			"Access-Control-Allow-Methods",
			"GET, PUT, POST, PATCH, DELETE",
		},
		{
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization",
		},
		{
			"Access-Control-Expose-Headers",
			"Authorization",
		},
		{
			"Access-Control-Max-Age",
			"600",
		},
	}

	for _, c := range cases {
		value := req.Header.Get(c.key)

		if value != c.expectedValue {
			t.Errorf("%v has not been set to the proper value, expect value to be %s, but get %s", c.key, c.expectedValue, value)
		}
	}

}
