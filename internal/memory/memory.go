package memory

import (
	"sync"
	"web-service/internal/models"
)

type MemoryStorage struct {
	Links map[string]*models.Link
	mu    sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Links: make(map[string]*models.Link),
		mu:    sync.RWMutex{},
	}
}

func (s *MemoryStorage) GetByShortCode(shortCode string, addVisits bool) (*models.Link, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, exists := s.Links[shortCode]
	if exists && addVisits {
		link.Visits++
	}
	return link, exists
}

func (s *MemoryStorage) DeleteByShortCode(shortCode string) *models.Link {
	s.mu.Lock()
	defer s.mu.Unlock()

	link, exists := s.Links[shortCode]
	if !exists {
		return nil
	}

	delete(s.Links, shortCode)

	return link
}
