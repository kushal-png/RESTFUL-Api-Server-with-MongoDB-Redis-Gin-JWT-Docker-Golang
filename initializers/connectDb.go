package initializers

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(config Config, ctx context.Context) (*mongo.Client, error) {
	connectionString := config.DBUri
	mongoconn := options.Client().ApplyURI(connectionString)
	mongoClient, err := mongo.Connect(ctx, mongoconn)
	return mongoClient, err
}

func ConnectRedis(config Config) *redis.Client {
	connectionString := config.RedisUri
	redisConn := redis.Options{Addr: connectionString}
	redisClient := redis.NewClient(&redisConn)
	return redisClient
}
