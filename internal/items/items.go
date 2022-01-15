package items

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
				log.Fatalf("Failed to JSONGet")
				return
			}
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
		todo.Done = false
		todo.Id = uuid.New()
		fmt.Println(todo)
		//DB
		res, err := server.Database.RedisJSONHandler.JSONSet("todo:"+todo.Id.String(), ".", todo)
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
	}
}

func RouterItems(server *structs.Server) {
	items := server.HttpRouter.Group("/items")
	{
		items.GET("/", getAllItems(server))
		items.POST("/", createItem(server))
		items.PATCH("/:todo-id/done", createItem(server))
	}
}
