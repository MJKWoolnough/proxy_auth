package proxy_auth

import (
	"encoding/json"
	"io"
)

var passFunc = NewPassHash256

// LoadSites takes json based data in an io.Reader
// JSON formatting must be of the following form
// [
//	{
//		"domain" : "domain1",
//		"users" : [
//			{ "username" : "user1", "password" : "password1" },
//			{ "username" : "user2", "password" : "password2" }
//		]
//	},
//	{
//		"domain" : "domain2",
//		"users" : [
//			{ "username" : "user3", "password" : "password3" },
//			{ "username" : "user4", "password" : "password4" }
//		]
//	}
// ]
func LoadSites(data io.Reader) (*Sites, error) {
	sd := []struct {
		Domain string `json:"domain"`
		Users  []struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"users"`
	}{}
	err := json.NewDecoder(data).Decode(&sd)
	if err != nil {
		return nil, err
	}

	s := NewSites()

	errs := make(ProcessErrors, 0)

	for _, site := range sd {
		u := NewSite()
		theseErrs := make(ProcessErrors, 0)
		for _, user := range site.Users {
			err = u.AddUser(user.Username, passFunc(user.Password))
			if err != nil {
				theseErrs = append(theseErrs, err)
			}
		}
		err = s.AddSite(site.Domain, u)
		if err != nil {
			errs = append(errs, err)
		} else {
			errs = append(errs, theseErrs...)
		}
	}
	if len(errs) > 0 {
		return s, errs
	}
	return s, nil
}

//Errors

type ProcessErrors []error

func (ProcessErrors) Error() string {
	return "errors encountered"
}
