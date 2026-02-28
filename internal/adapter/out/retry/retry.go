package retry

import (
	"time"

	"github.com/wb-go/wbf/retry"
)

const (
	DefaultAttempts  = 3
	DefaultDelay     = 500 * time.Millisecond
	DefaultBackoff   = 2.0
	DatabaseAttempts = 5
	DatabaseDelay    = 1 * time.Second
	DatabaseBackoff  = 2.0
	RedisAttempts    = 3
	RedisDelay       = 300 * time.Millisecond
	RedisBackoff     = 1.5
	ServiceAttempts  = 3
	ServiceDelay     = 200 * time.Millisecond
	ServiceBackoff   = 1.5
)

func GetDatabaseStrategy() retry.Strategy {
	return retry.Strategy{
		Attempts: DatabaseAttempts,
		Delay:    DatabaseDelay,
		Backoff:  DatabaseBackoff,
	}
}

func GetRedisStrategy() retry.Strategy {
	return retry.Strategy{
		Attempts: RedisAttempts,
		Delay:    RedisDelay,
		Backoff:  RedisBackoff,
	}
}

func GetServiceStrategy() retry.Strategy {
	return retry.Strategy{
		Attempts: ServiceAttempts,
		Delay:    ServiceDelay,
		Backoff:  ServiceBackoff,
	}
}
