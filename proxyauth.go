// proxy_auth implements the Top Coder - Go Learning Challenge - Simple Web-API Server
package proxy_auth

import (
	"io"
	"net/http"
)

var (
	APIPath = "/api/2/domains/"
	APIFile = "/proxyauth"
)

func Setup(data io.Reader, sm *http.ServeMux) error {
	s, err := LoadSites(data)
	if err != nil {
		return err
	}
	sm.Handle(APIPath, NewServer(APIPath, APIFile, s))
	return nil
}
