package proto

import (
	"context"
)

type PingService interface {
	Ping(ctx context.Context) error
}

type StatService interface {
	GetUrlsAndUsersStat() (int32, int32)
}

type Server struct {
	UnimplementedServiceServer
	pingService PingService
	statService StatService
}

func NewGPRCServer(pingService PingService, statService StatService) *Server {
	return &Server{
		pingService: pingService,
		statService: statService,
	}

}

func (s *Server) GetStat(ctx context.Context, r *StatRequest) (*StatsResponse, error) {
	var response StatsResponse

	urlCount, userCount := s.statService.GetUrlsAndUsersStat()

	response.Urls = urlCount
	response.Users = userCount

	return &response, nil
}

func (s *Server) Ping(ctx context.Context, r *PingRequest) (*PingResponse, error) {
	var response PingResponse

	err := s.pingService.Ping(ctx)
	if err != nil {
		response.Error = err.Error()
	}

	return &response, nil
}
