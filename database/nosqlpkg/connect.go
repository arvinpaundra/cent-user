package nosqlpkg

import (
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

type InMemoryConnectible interface {
	Connect() (*redis.Client, error)
	Close() error
}

func NewInMemoryConection(connect InMemoryConnectible) {
	once.Do(func() {
		var err error

		rdb, err = connect.Connect()
		if err != nil {
			log.Fatalf("failed to in memory database: %s", err.Error())
		}
	})
}

func GetInMemoryConnection() *redis.Client {
	return rdb
}
