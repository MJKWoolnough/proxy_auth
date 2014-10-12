package proxy_auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewSites()

	examplecom := NewSite()
	testorg := NewSite()
	hellonet := NewSite()

	examplecom.AddUser("user1", NewPassHash256("password1"))
	examplecom.AddUser("user2", NewPassHash256("abcd1234"))
	examplecom.AddUser("myname", NewPassHash256("drowssap"))
	examplecom.AddUser("joe", NewPassHash256("letmein"))
	examplecom.AddUser("steve", NewPassHash256("enter"))

	testorg.AddUser("user1", NewPassHash256("password1"))
	testorg.AddUser("hello", NewPassHash256("world"))

	hellonet.AddUser("foo", NewPassHash256("bar"))

	s.AddSite("example.com", examplecom)
	s.AddSite("test.org", testorg)
	s.AddSite("hello.net", hellonet)

	server := NewServer("/", "/api", s)

	tests := []struct {
		domain, username, passhash string
		responseCode               int
		responseBody               []byte
	}{
		{"example.com", "user1", "{SHA256}CxTVAaWURCoBxoWVQbyz6BZNGD0yk3uFGDVEL2nVyU4=", 200, Access_Granted},
		{"example.com", "user2", "{SHA256}CxTVAaWURCoBxoWVQbyz6BZNGD0yk3uFGDVEL2nVyU4=", 200, Access_Denied},
		{"test.org", "user1", "{SHA256}CxTVAaWURCoBxoWVQbyz6BZNGD0yk3uFGDVEL2nVyU4=", 200, Access_Granted},
		{"test.org", "hello", "{SHA256}SG6kYiTRu0+2gPNPfJrZao8k7Ii+c+qOWmxlJg6cuKc=", 200, Access_Granted},
		{"hello.net", "foo", "{SHA256}/N4rLtula/QIYB+3If6bXDONEO5CnqBPrlURto+/j7k=", 200, Access_Granted},
		{"hello.com", "foo", "{SHA256}/N4rLtula/QIYB+3If6bXDONEO5CnqBPrlURto+/j7k=", 404, []byte{}},
	}

	for n, test := range tests {
		w := httptest.NewRecorder()
		v := url.Values{"username": {test.username}, "password": {test.passhash}}
		r, _ := http.NewRequest("POST", "http://localhost/"+test.domain+"/api", strings.NewReader(v.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		server.ServeHTTP(w, r)
		if w.Code != test.responseCode {
			t.Errorf("test %d: expecting response code %d, got %d", n+1, test.responseCode, w.Code)
		} else if !bytes.Equal(w.Body.Bytes(), test.responseBody) {
			t.Errorf("test %d: expecting body %s, got %s", n+1, w.Body, string(test.responseBody))
		}
	}
}
