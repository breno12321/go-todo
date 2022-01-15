package structs

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type Database struct {
	RedisJSONHandler *rejson.Handler
	RedisClient      *redis.Client
}

type Server struct {
	Database   *Database
	HttpRouter *gin.RouterGroup
	Context    context.Context
}
