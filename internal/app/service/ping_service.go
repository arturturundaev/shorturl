package service

type PingService struct {
	repository PingRepository
}

func NewPingService(repository PingRepository) *PingService {
	return &PingService{repository: repository}
}

func (s *PingService) Ping() error {
	return s.repository.Ping()
}
