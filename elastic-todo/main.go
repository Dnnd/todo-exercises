package main

import (
	"context"
	"github.com/Dnnd/todo-exercises/elastic-todo/api"
	"github.com/Dnnd/todo-exercises/elastic-todo/models"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
)

func Map(client *elastic.Client, model models.ElasticModel) {
	ctx := context.Background()
	if exists, _ := client.IndexExists("user").Do(ctx); exists == false {
		if _, err := client.CreateIndex("user").BodyString(model.GetMapping()).Do(ctx); err != nil {
			panic(err)
		}
	}
}

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		panic(err)
	}
	Map(client, models.User{})
	Map(client, models.Todo{})

	defer client.Stop()
	r := gin.Default()
	service := api.Service{
		Elastic: client,
	}
	r.POST("/users/:user_id/todos", service.MakeTodo)
	r.GET("/users/:user_id/todos/:todo_id", service.FetchTodo)
	r.GET("/users/:user_id/todos/", service.FetchUserTodos)
	r.PUT("/users/:user_id/todos/:todo_id", service.UpdateTodo)
	r.DELETE("/users/:user_id/todos/:todo_id", service.DeleteTodo)
	r.GET("/users/:user_id", service.GetUser)
	r.GET("/users", service.FilterUsers)
	r.POST("/users", service.MakeUser)
	r.Run("localhost:10000")
}
