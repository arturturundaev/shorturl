package service

import "context"

// PingService сервис
type PingService struct {
	repository RepositoryReader
}

// NewPingService Конструктор
func NewPingService(repository RepositoryReader) *PingService {
	return &PingService{repository: repository}
}

// Ping пиннг
func (s *PingService) Ping(ctx context.Context) error {
	return s.repository.Ping(ctx)
}
