package follows

import "sync"

type Service struct {
	mu        sync.RWMutex
	following map[string]map[string]struct{}
	followers map[string]map[string]struct{}
}

func NewService() *Service {
	return &Service{
		following: make(map[string]map[string]struct{}),
		followers: make(map[string]map[string]struct{}),
	}
}

func (s *Service) Follow(followerID, followeeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.following[followerID]; !ok {
		s.following[followerID] = make(map[string]struct{})
	}
	if _, ok := s.followers[followeeID]; !ok {
		s.followers[followeeID] = make(map[string]struct{})
	}
	s.following[followerID][followeeID] = struct{}{}
	s.followers[followeeID][followerID] = struct{}{}
}

func (s *Service) Unfollow(followerID, followeeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if items, ok := s.following[followerID]; ok {
		delete(items, followeeID)
	}
	if items, ok := s.followers[followeeID]; ok {
		delete(items, followerID)
	}
}

func (s *Service) CountFollowing(userID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.following[userID])
}

func (s *Service) CountFollowers(userID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.followers[userID])
}
