package storage

type Config interface {
	ReturnDatabase() string
	ReturnUser() string
	ReturnPassword() string
	ReturnHost() string
	ReturnPort() int
	ReturnRetries() int
	ReturnPoolSize() int
}
