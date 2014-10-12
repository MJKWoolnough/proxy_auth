package proxy_auth

import (
	"bytes"
	"reflect"
	"testing"
)

func TestLoadSiteS(t *testing.T) {
	tests := []struct {
		json string
		errs error
	}{
		{
			`[ { "domain": "test1.com", "users": [ { "username": "user1", "password": "pass1" } ] } ]`,
			nil,
		},
		{
			`[ { "domain": "samesite.net", "users" : [] }, { "domain": "samesite.net", "users" : [] } ]`,
			ProcessErrors{SiteExists("samesite.net")},
		},
		{
			`[ { "domain": "test1.com", "users": [ { "username": "sameuser", "password": "pass1" },` +
				`{ "username": "sameuser", "password": "pass2" } ] } ]`,
			ProcessErrors{UserExists("sameuser")},
		},
		{
			`[` +
				`	{` +
				`		"domain": "testme.co.jp",` +
				`		"users": [` +
				`			{ "username": "user1", "password": "password1" },` +
				`			{ "username": "user2", "password": "mypassword" }` +
				`		]` +
				`	},` +
				`	{` +
				`		"domain": "example.com",` +
				`		"users": [` +
				`			{ "username": "user1", "password": "password1" },` +
				`			{ "username": "user2", "password": "mypassword" }` +
				`		]` +
				`	},` +
				`	{` +
				`		"domain": "testme.co.jp", ` + //duplicate site
				`		"users": [` +
				`			{ "username": "sameuser", "password": "password1" },` + //Duplicate users on duplicate domains
				`			{ "username": "sameuser", "password": "mypassword" }` + //should not cause user based errors
				`		]` +
				`	},` +
				`	{` +
				`		"domain": "example.org",` +
				`		"users": [` +
				`			{ "username": "user", "password": "password1" },` + //duplicate user
				`			{ "username": "user", "password": "mypassword" }` +
				`		]` +
				`	}` +
				`]`,
			ProcessErrors{SiteExists("testme.co.jp"), UserExists("user")},
		},
	}

	for n, test := range tests {
		_, errs := LoadSites(bytes.NewReader([]byte(test.json)))
		if !reflect.DeepEqual(errs, test.errs) {
			t.Errorf("test %d: received errors did not match expected", n+1)
		}
	}
}
