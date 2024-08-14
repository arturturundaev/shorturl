package service

type RepositoryPinger interface {
	Ping() error
}
