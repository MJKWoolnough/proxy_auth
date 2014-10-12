package proxy_auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
)

type Password interface {
	Auth(string) error
	Set(string) error
}

type PassHash256 [sha256.Size]byte

const passHash256Prefix = "{SHA256}"
const passHash245Len = 52

func NewPassHash256(password string) Password {
	p := new(PassHash256)
	p.Set(password)
	return p
}

func (p *PassHash256) Set(password string) error {
	*p = PassHash256(sha256.Sum256([]byte(password)))
	return nil
}

// Auth takes in a base64 encoded sha256 hash with "{SHA256}" (without quotes" prepended)
// and matches against the stored password.
//
// A nil error represents a successful authentication.
func (p *PassHash256) Auth(password string) error {
	if len(password) != passHash245Len {
		return IncorrectPasswordLength(len(password))
	}
	if password[:len(passHash256Prefix)] != passHash256Prefix {
		return IncorrectPasswordPrefix{}
	}
	if q, err := base64.StdEncoding.DecodeString(password[len(passHash256Prefix):]); err != nil {
		return PasswordParseError{err}
	} else if !bytes.Equal((*p)[:], q) { //bytes.Equal uses bytes·Equal ASM for effeciency
		return IncorrectPassword{}
	}
	return nil

}

//Errors

type IncorrectPasswordLength int

func (i IncorrectPasswordLength) Error() string {
	return "incorrect password length, expecting " + strconv.Itoa(passHash245Len) + ", got " + strconv.Itoa(int(i))
}

type IncorrectPasswordPrefix struct{}

func (IncorrectPasswordPrefix) Error() string {
	return "incorrect password prefix"
}

type PasswordParseError struct {
	Err error
}

func (p PasswordParseError) Error() string {
	return "error occurred while parsing password: " + p.Err.Error()
}

type IncorrectPassword struct{}

func (IncorrectPassword) Error() string {
	return "incorrect password"
}
