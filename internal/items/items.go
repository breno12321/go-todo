package items

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/breno12321/go-todo/internal/helpers"
	"github.com/breno12321/go-todo/internal/structs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Todo struct {
	Id          uuid.UUID `json:"id,omitempty"`
	Title       string    `json:"title,omitempty" binding:"required"`
	Description string    `json:"description,omitempty" binding:"required"`
	Done        bool      `json:"done,omitempty"`
	CreatedAt   time.Time `json:"created-at,omitempty"`
	UpdatedAt   time.Time `json:"updated-at,omitempty"`
}

func getAllItems(server *structs.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		todosKeys, err := server.Database.RedisClient.Keys(server.Context, "todo:*").Result()
		if err != nil {
			log.Fatal("Failed to get Keys")
			return
		}

		allTodo := []Todo{}
		for _, v := range todosKeys {
			todoJSON, err := server.Database.RedisJSONHandler.JSONGet(v, ".")
			if err != nil {
				log.Fatalf("Failed to JSONGet %v", err)
				return
			}
			todoEntry := Todo{}
			err = json.Unmarshal(todoJSON.([]byte), &todoEntry)
			if err != nil {
				log.Fatalf("Failed to JSON Unmarshal %v", err)
				return
			}
			allTodo = append(allTodo, todoEntry)
		}

		// fmt.Printf("Student read from redis : %#v\n", readTodo)
		c.JSON(http.StatusOK, gin.H{"data": allTodo})
	}
}

func getItem(server *structs.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		todoId := c.Param("todo-id")

		todoJSON, err := server.Database.RedisJSONHandler.JSONGet("todo:"+todoId, ".")
		if err != nil {
			helpers.ErrorLogger.Printf("Failed to JSONGet error %v", err)
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		todoEntry := Todo{}
		err = json.Unmarshal(todoJSON.([]byte), &todoEntry)
		if err != nil {
			helpers.ErrorLogger.Printf("Failed to JSON Unmarshal %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		// fmt.Printf("Student read from redis : %#v\n", readTodo)
		c.JSON(http.StatusOK, gin.H{"data": todoEntry})
	}
}

func createItem(server *structs.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		//CONTROLLER PARSE
		var todo Todo
		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//ENTITY CREATION
		todo.CreatedAt = time.Now()
		todo.UpdatedAt = time.Now()
		todo.Id = uuid.New()
		todo.Done = true
		fmt.Println(todo)

		//DB
		res, err := server.Database.RedisJSONHandler.JSONSet("todo:"+todo.Id.String(), ".", todo)
		if err != nil {
			log.Fatalf("Failed to JSONSet")
			return
		}
		fmt.Println(res)

		//RESPONSE
		c.JSON(http.StatusCreated, gin.H{"data": todo})
	}
}

func toggleItemDone(server *structs.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		//CONTROLLER PARSE
		todoId := c.Param("todo-id")

		//ENTITY CREATION
		toggleAt := time.Now()

		//DB

		// Store updated value
		_, err := server.Database.RedisClient.Do(server.Context, "JSON.TOGGLE", "todo:"+todoId, ".done").Result()
		if err != nil {
			helpers.ErrorLogger.Printf("Failed to JSONSet error %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		res, err := server.Database.RedisJSONHandler.JSONSet("todo:"+todoId, ".updated-at", toggleAt)
		if err != nil {
			helpers.ErrorLogger.Printf("Failed to JSONSet error %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if res.(string) == "OK" {
			fmt.Printf("Success: %s\n", res)
		} else {
			fmt.Println("Failed to Set: ")
		}

		helpers.InfoLogger.Printf("Item %v toggled", todoId)

		//RESPONSE
		c.JSON(http.StatusOK, gin.H{})
	}
}

func RouterItems(server *structs.Server) {
	items := server.HttpRouter.Group("/items")
	{
		items.GET("/", getAllItems(server))
		items.GET("/:todo-id", getItem(server))
		items.POST("/", createItem(server))
		items.POST("/:todo-id/toggle-done", toggleItemDone(server))
	}
}
