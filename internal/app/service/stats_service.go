package service

type StatService struct {
	repository getStatsRepository
}

type getStatsRepository interface {
	GetUrlsCount() int32
	GetUsersCount() int32
}

func NewStatService(repository getStatsRepository) *StatService {
	return &StatService{repository: repository}
}

func (s *StatService) GetUrlsAndUsersStat() (int32, int32) {
	urlsCount := s.repository.GetUrlsCount()
	usersCount := s.repository.GetUsersCount()

	return urlsCount, usersCount
}
