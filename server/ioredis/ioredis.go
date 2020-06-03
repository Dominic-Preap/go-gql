package ioredis

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis/v7"
	"github.com/my/app/server/config"
)

// InitRedis Create and connect to redis server
func InitRedis(env *config.EnvConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     env.RedisAddress,
		Password: env.RedisPassword,
		DB:       0, // use default DB,
	})

	pong, err := client.Ping().Result()
	log.Printf("Redis Connected: %v, %v", pong, err)

	// For Testing Expiry
	// client.Set("ex:reminder:1", "test", time.Duration(time.Second*5))
	// client.Set("ex:notification:2", "test", time.Duration(time.Second*10))

	return client
}

// InitRedisExpiryPubSub Create as pub/sub service.
func InitRedisExpiryPubSub(s *config.Server) {
	// Keyspace notifications allow clients to subscribe to Pub/Sub channels
	// in order to receive events affecting the Redis data set in some way.
	s.Client.ConfigSet("notify-keyspace-events", "Ex")

	expired := make(chan redis.Message)
	go func() {
		redisChan := s.Client.PSubscribe("__keyevent@0__:expired").Channel()
		for redisMsg := range redisChan {
			expired <- *redisMsg
		}
	}()

	// handle functions on redis expiry events
	go func() {
		for msg := range expired {
			// log.Printf("debug: expired %+v", msg)

			// ! * Naming Convention : ex:TYPE:KEY
			split := strings.Split(msg.Payload, ":")
			key := split[1]
			value := split[2]

			switch key {
			case "reminder":
				reminder(s, value)
			case "notification":
				notification(s, value)
			default:
				break
			}

		}
	}()
}

func reminder(s *config.Server, value string) {
	// Your logic here
	fmt.Printf("Ex Reminder: %v \n", value)
}

func notification(s *config.Server, value string) {
	// Your logic here
	fmt.Printf("Ex Notification: %v \n", value)
}
