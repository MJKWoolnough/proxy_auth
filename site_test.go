package proxy_auth

import (
	"reflect"
	"testing"
)

const (
	AddUser CRUDFunc = iota
	UpdateUser
	RemoveUser
	Auth
)

type CRUDFunc int

func TestSiteUserCRUD(t *testing.T) {
	t.Parallel()

	s := NewSite()

	tests := []struct {
		username, password string
		cf                 CRUDFunc
		err                error
		passwordType       Password
	}{
		{"testUser", "myPassword", AddUser, nil, new(PassHash256)},
		{"testUser", "myPassword", AddUser, UserExists("testUser"), new(PassHash256)},
		{"testUser2", "myPassword", AddUser, nil, new(PassHash256)},
		{"testUser", "myNewPassword", UpdateUser, nil, new(PassHash256)},
		{"testUser3", "myNewPassword", UpdateUser, NoUser("testUser3"), new(PassHash256)},
		{"testUser2", "myNewPassword", RemoveUser, nil, new(PassHash256)},
		{"testUser3", "myNewPassword", RemoveUser, NoUser("testUser3"), new(PassHash256)},
		{"testUser3", "myNewPassword", Auth, NoUser("testUser3"), new(PassHash256)},
		{"testUser2", "myPassword2", AddUser, nil, new(PassHash256)},
	}

	for n, test := range tests {
		var err error
		test.passwordType.Set(test.password)
		switch test.cf {
		case AddUser:
			err = s.AddUser(test.username, test.passwordType)
		case UpdateUser:
			err = s.UpdateUser(test.username, test.passwordType)
		case RemoveUser:
			err = s.RemoveUser(test.username)
		case Auth:
			err = s.Auth(test.username, "")
		}
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("test %d: expecting error %q, got %q", n+1, test.err, err)
		}
	}
}
