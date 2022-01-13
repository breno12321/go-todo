package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/nitishm/go-rejson/v4"
)

var ctx = context.Background()

// func getAllItems(rh *rejson.Handler) []Todo {

// }

type Todo struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Title       string    `json:"title,omitempty" binding:"required"`
	Description string    `json:"description,omitempty" binding:"required"`
	Done        bool      `json:"done,omitempty"`
	CreatedAt   time.Time `json:"created-at,omitempty"`
}

type Database struct {
	RedisJSONHandler *rejson.Handler
	RedisClient      *redis.Client
}

type Server struct {
	Database   Database
	HttpRouter gin.RouterGroup
}

func routerItems(r *gin.RouterGroup, db *Database) {
	items := r.Group("/items")
	{
		items.GET("/", func(c *gin.Context) {

			todosKeys, err := db.RedisClient.Keys(ctx, "todo:*").Result()
			if err != nil {
				log.Fatal("Failed to get Keys")
				return
			}

			allTodo := []Todo{}
			for _, v := range todosKeys {
				todoJSON, err := db.RedisJSONHandler.JSONGet(v, ".")
				if err != nil {
					log.Fatalf("Failed to JSONGet")
					return
				}
				todoas, err := db.RedisClient.Do(ctx, "JSONGET", v, ".").Result()
				fmt.Println(string(todoJSON.([]byte)))
				fmt.Println(todoas)
				todoEntry := Todo{}
				err = json.Unmarshal(todoJSON.([]byte), &todoEntry)
				if err != nil {
					log.Fatalf("Failed to JSON Unmarshal")
					return
				}
				allTodo = append(allTodo, todoEntry)
			}

			// fmt.Printf("Student read from redis : %#v\n", readTodo)
			c.JSON(http.StatusOK, gin.H{"items": allTodo})
		})
		items.POST("/", func(c *gin.Context) {
			//CONTROLLER PARSE
			var todo Todo
			if err := c.ShouldBindJSON(&todo); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			//ENTITY CREATION
			todo.CreatedAt = time.Now()
			todo.Done = false
			todo.Id = uuid.New()

			//DB
			res, err := db.RedisJSONHandler.JSONSet("todo:"+todo.Id.String(), ".", todo)
			if err != nil {
				log.Fatalf("Failed to JSONSet")
				return
			}
			if res.(string) == "OK" {
				fmt.Printf("Success: %s\n", res)
			} else {
				fmt.Println("Failed to Set: ")
			}

			//RESPONSE
			c.JSON(http.StatusCreated, gin.H{"todo": todo})
		})
	}
}

func main() {
	router := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClient(rdb)

	db := Database{RedisJSONHandler: rh,
		RedisClient: rdb}

	//API Routes
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})
		routerItems(api, &db)
	}

	router.Run(":8080")
}
