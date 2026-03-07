package main

import "sync"

type MemStore struct {
	mu   sync.RWMutex
	data map[string]QRCode
}

func NewMemStore() *MemStore {
	return &MemStore{data: make(map[string]QRCode)}
}

func (s *MemStore) Create(qr QRCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[qr.ID] = qr
	return nil
}

func (s *MemStore) List() ([]QRCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]QRCode, 0, len(s.data))
	for _, qr := range s.data {
		list = append(list, qr)
	}
	return list, nil
}

func (s *MemStore) GetByID(id string) (QRCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	qr, ok := s.data[id]
	if !ok {
		return QRCode{}, ErrNotFound
	}
	return qr, nil
}

func (s *MemStore) GetByClientNumber(phone string) (QRCode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, qr := range s.data {
		if qr.ClientNumber == phone {
			return qr, nil
		}
	}
	return QRCode{}, ErrNotFound
}
