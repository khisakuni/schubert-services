package service

type Auth interface {
	HashPassword(string) (string, error)
}
