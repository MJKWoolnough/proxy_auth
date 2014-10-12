package proxy_auth

import (
	"bytes"
	"encoding/base64"
	"reflect"
	"testing"
)

func TestNewPassHash256(t *testing.T) {
	t.Parallel()
	tests := []struct {
		password string
		hash     []byte
	}{
		{"abc123", []byte{
			108, 161, 61, 82, 202, 112, 200, 131, 224, 240, 187, 16, 30, 66, 90, 137,
			232, 98, 77, 229, 29, 178, 210, 57, 37, 147, 175, 106, 132, 17, 128, 144,
		}},
		{"abcd1234", []byte{
			233, 206, 231, 26, 185, 50, 253, 232, 99, 51, 141, 8, 190, 77, 233, 223,
			227, 158, 160, 73, 189, 175, 179, 66, 206, 101, 158, 197, 69, 11, 105, 174,
		}},
		{"ilovego", []byte{
			217, 2, 112, 111, 77, 34, 200, 214, 153, 110, 193, 27, 141, 129, 212, 77,
			50, 242, 189, 28, 36, 37, 148, 237, 243, 42, 227, 226, 161, 214, 5, 53,
		}},
	}
	p := new(PassHash256)
	for n, test := range tests {
		p.Set(test.password)
		if !bytes.Equal(p[:], test.hash) {
			t.Errorf("test %d: expecting hash %v, got %v", n+1, test.hash, p[:])
		}
	}
}

func TestPassHash256Auth(t *testing.T) {
	t.Parallel()

	p := new(PassHash256)

	tests := []struct {
		password, apiString string
		err                 error
	}{
		{"password", "", IncorrectPasswordLength(0)},
		{"abc123", "{SHA256}bKE9UspwyIPg8LsQHkJaiehiTeUdstI5JZOvaoQRgJA", IncorrectPasswordLength(51)},                        //missing trailing =
		{"abcd1234", "{SHA255}6c7nGrky/ehjM40Ivk3p3+OeoEm9r7NCzmWexUULaa4=", IncorrectPasswordPrefix{}},                       //255, not 256
		{"ilovego", "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBT*=", PasswordParseError{base64.CorruptInputError(42)}}, //invalid base64
		{"wrongPassword", "{SHA256}56fDF6gfgD4sdRfg5SDTbKsSD566gDDvbFFFDH73FG2=", IncorrectPassword{}},                        //incorrect
		{"abc123", "{SHA256}bKE9UspwyIPg8LsQHkJaiehiTeUdstI5JZOvaoQRgJA=", nil},                                               //correct
		{"abcd1234", "{SHA256}6c7nGrky/ehjM40Ivk3p3+OeoEm9r7NCzmWexUULaa4=", nil},                                             //correct
		{"ilovego", "{SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=", nil},                                              //correct
	}

	for n, test := range tests {
		p.Set(test.password)
		if err := p.Auth(test.apiString); !reflect.DeepEqual(err, test.err) {
			if test.err == nil {
				t.Errorf("test %d: received unexpected error: %q", n+1, err)
			} else {
				t.Errorf("test %d: expecting error %q, got %q", n+1, test.err, err)
			}
		}
	}
}
