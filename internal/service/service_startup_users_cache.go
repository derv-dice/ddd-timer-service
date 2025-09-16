package service

import "context"

func (s *Service) StartupUsersCache(ctx context.Context) error {
	users, err := s.repo.LoadAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		s.usersCache.Set(user.ID, user)
	}

	return nil
}
