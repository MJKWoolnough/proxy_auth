package proxy_auth

import (
	"net/http"
	"strings"
)

var (
	Access_Granted = []byte("{\n	\"access_granted\": true\n}")
	Access_Denied  = []byte("{\n	\"access_granted\": false, \"reason\": \"denied by policy\"\n}")
)

type Server struct {
	prefix, suffix string
	*Sites
}

func NewServer(prefix, suffix string, s *Sites) *Server {
	return &Server{
		prefix: prefix,
		suffix: suffix,
		Sites:  s,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) <= len(s.prefix)+len(s.suffix) ||
		r.URL.Path[:len(s.prefix)] != s.prefix ||
		r.URL.Path[len(r.URL.Path)-len(s.suffix):] != s.suffix {
		w.WriteHeader(http.StatusNotFound)
	} else {
		domain := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, s.prefix), s.suffix)
		r.ParseForm()
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		err := s.Auth(domain, username, password)
		if err == nil {
			s.writeResult(w, Access_Granted)
		} else {
			switch err.(type) {
			case NoSite:
				w.WriteHeader(http.StatusNotFound)
			default:
				s.writeResult(w, Access_Denied)
			}
		}
	}
}

func (s *Server) writeResult(w http.ResponseWriter, data []byte) {
	w.Header()["Content-Type"] = []string{"application/json"}
	w.Write(data)
}
