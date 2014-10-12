package proxy_auth

import (
	"reflect"
	"testing"
)

const (
	AddSite CRUDFunc = iota
	RemoveSite
	GetSite
)

func TestSitesCRUD(t *testing.T) {
	t.Parallel()

	s := NewSites()

	tests := []struct {
		domain string
		cf     CRUDFunc
		err    error
	}{
		{"example.com", AddSite, nil},
		{"test.net", AddSite, nil},
		{"example.com", AddSite, SiteExists("example.com")},
		{"example.com", GetSite, nil},
		{"nothere.org", GetSite, NoSite("nothere.org")},
		{"test.net", RemoveSite, nil},
		{"nothere2.org", RemoveSite, NoSite("nothere2.org")},
		{"test.net", AddSite, nil},
	}

	for n, test := range tests {
		var err error
		switch test.cf {
		case AddSite:
			err = s.AddSite(test.domain, NewSite())
		case RemoveSite:
			err = s.RemoveSite(test.domain)
		case GetSite:
			_, err = s.getSite(test.domain)
		}
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("test %d: expecting error %q, got %q", n+1, test.err, err)
		}
	}
}
