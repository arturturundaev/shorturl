package service

import (
	"context"
)

type PingService struct {
	repository PingRepository
}

func NewPingService(repository PingRepository) *PingService {
	return &PingService{repository: repository}
}

func (s *PingService) Ping(ctx context.Context) error {
	return s.repository.Ping(ctx)
}
