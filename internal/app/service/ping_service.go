package service

type PingService struct {
	repository RepositoryReadInterface
}

func NewPingService(repository RepositoryReadInterface) *PingService {
	return &PingService{repository: repository}
}

func (s *PingService) Ping() error {
	return s.repository.Ping()
}
