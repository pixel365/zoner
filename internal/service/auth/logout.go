package auth

func (s *Service) Logout() error {
	return s.repo.Logout()
}
