package proxy_auth

import "sync"

type Sites struct {
	sites map[string]*Site
	l     sync.RWMutex
}

func NewSites() *Sites {
	return &Sites{
		sites: make(map[string]*Site),
	}
}

func (s *Sites) AddSite(domain string, site *Site) error {
	s.l.Lock()
	defer s.l.Unlock()
	if _, ok := s.sites[domain]; ok {
		return SiteExists(domain)
	}
	s.sites[domain] = site
	return nil
}

func (s *Sites) RemoveSite(domain string) error {
	s.l.Lock()
	defer s.l.Unlock()
	if _, ok := s.sites[domain]; !ok {
		return NoSite(domain)
	}
	delete(s.sites, domain)
	return nil
}

func (s *Sites) getSite(domain string) (*Site, error) {
	s.l.RLock()
	defer s.l.RUnlock()
	if site, ok := s.sites[domain]; ok {
		return site, nil
	} else {
		return nil, NoSite(domain)
	}
}

func (s *Sites) AddUser(domain, user string, password Password) error {
	site, err := s.getSite(domain)
	if err != nil {
		return err
	}
	return site.AddUser(user, password)
}

func (s *Sites) UpdateUser(domain, user string, password Password) error {
	site, err := s.getSite(domain)
	if err != nil {
		return err
	}
	return site.UpdateUser(user, password)
}

func (s *Sites) RemoveUser(domain, user string) error {
	site, err := s.getSite(domain)
	if err != nil {
		return err
	}
	return site.RemoveUser(user)
}

func (s *Sites) Auth(domain, user, password string) error {
	site, err := s.getSite(domain)
	if err != nil {
		return err
	}
	return site.Auth(user, password)
}

//Errors

type NoSite string

func (n NoSite) Error() string {
	return "site " + string(n) + " does not exist"
}

type SiteExists string

func (s SiteExists) Error() string {
	return "cannot add side " + string(s) + " - already exists"
}
