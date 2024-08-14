package service

import "context"

type PingService struct {
	repository RepositoryReader
}

func NewPingService(repository RepositoryReader) *PingService {
	return &PingService{repository: repository}
}

func (s *PingService) Ping(ctx context.Context) error {
	return s.repository.Ping(ctx)
}
