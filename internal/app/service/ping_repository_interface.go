package service

type PingRepository interface {
	Ping() error
}
