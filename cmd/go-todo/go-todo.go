package main

import (
	"context"
	"net/http"

	"github.com/breno12321/go-todo/internal/items"
	"github.com/breno12321/go-todo/internal/structs"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

var ctx = context.Background()

type Database = structs.Database

type Server = structs.Server

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClient(rdb)

	db := &structs.Database{RedisJSONHandler: rh,
		RedisClient: rdb}

	//API Routes
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})
		// items.RouterItems(api, db, ctx)
		items.RouterItems(&structs.Server{Database: db, HttpRouter: api, Context: ctx})
	}

	router.Run(":8080")
}
