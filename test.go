// +build ignore

package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	u        = flag.String("url", "", "url including {domain}")
	userfile = flag.String("testfile", "", "json test file")
)

func password2PassHash256(password string) string {
	hash := sha256.Sum256([]byte(password))
	return "{SHA256}" + base64.StdEncoding.EncodeToString(hash[:])
}

type userData struct {
	Domain, Username, Password string
	ExpectingCode              int
	ShouldAuth                 bool
}

type responseData struct {
	AccessGranted bool   `json:"access_granted"`
	Reason        string `json:"reason"`
}

func main() {
	flag.Parse()
	if *userfile != "" {
		f, err := os.Open(*userfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := make([]userData, 0)
		err = json.NewDecoder(f).Decode(&data)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, user := range data {
			fmt.Printf("Testing user %q @ %q, with password %q\n", user.Username, user.Domain, user.Password)
			user.Domain = strings.Replace(*u, "{domain}", user.Domain, 1)
			err = TestUser(user)
			if err != nil {
				fmt.Printf("\x1b[31mFAILED\x1b[39m\n	===ERROR=== %s\n", err)
			} else {
				fmt.Println("\x1b[32mPASSED\x1b[39m")
			}
		}
	}
}

func TestUser(user userData) error {
	r, err := http.PostForm(user.Domain, url.Values{"username": {user.Username}, "password": {password2PassHash256(user.Password)}})
	if err != nil {
		return err
	}
	if r.StatusCode != user.ExpectingCode {
		return StatusError{r.StatusCode, user.ExpectingCode}
	}
	rd := new(responseData)
	err = json.NewDecoder(r.Body).Decode(rd)
	if err != nil && err != io.EOF {
		return err
	}
	if rd.AccessGranted != user.ShouldAuth {
		return AccessError{rd.AccessGranted, user.ShouldAuth}
	}
	return nil
}

//Errors

type StatusError struct {
	got, expecting int
}

func (s StatusError) Error() string {
	return fmt.Sprintf("expecting status code %d, got %d", s.expecting, s.got)
}

type AccessError struct {
	got, expecting bool
}

func (a AccessError) Error() string {
	return fmt.Sprintf("expecting access_granted = %v, got %v", a.expecting, a.got)
}
