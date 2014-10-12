package proxy_auth

import "sync"

type Site struct {
	users map[string]Password
	l     sync.RWMutex
}

func NewSite() *Site {
	return &Site{
		users: make(map[string]Password),
	}
}

func (s *Site) AddUser(username string, password Password) error {
	s.l.Lock()
	defer s.l.Unlock()
	if _, ok := s.users[username]; ok {
		return UserExists(username)
	}
	s.users[username] = password
	return nil
}

func (s *Site) UpdateUser(username string, password Password) error {
	s.l.Lock()
	defer s.l.Unlock()
	if _, ok := s.users[username]; !ok {
		return NoUser(username)
	}
	s.users[username] = password
	return nil

}

func (s *Site) RemoveUser(username string) error {
	s.l.Lock()
	defer s.l.Unlock()
	if _, ok := s.users[username]; !ok {
		return NoUser(username)
	}
	delete(s.users, username)
	return nil
}

func (s *Site) Auth(username, password string) error {
	s.l.RLock()
	defer s.l.RUnlock()
	if u, ok := s.users[username]; ok {
		return u.Auth(password)
	}
	return NoUser(username)
}

//Errors

type NoUser string

func (n NoUser) Error() string {
	return "user " + string(n) + " not found"
}

type UserExists string

func (u UserExists) Error() string {
	return "cannot add user " + string(u) + " - already exists"
}
