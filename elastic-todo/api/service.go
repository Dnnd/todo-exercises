package api

import (
	"context"
	"github.com/Dnnd/todo-exercises/elastic-todo/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"log"
	"net/http"
)

type Service struct {
	Elastic *elastic.Client
}

func (s *Service) MakeTodo(ctx *gin.Context) {
	client := s.Elastic
	reqCtx := context.Background()
	var todo models.Todo
	ctx.BindJSON(&todo)
	userId := ctx.Param("user_id")
	res, err := client.Exists().
		Index("user").
		Type("_doc").
		Id(userId).
		Do(reqCtx)
	if  !res {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	todo.UserId = userId
	created, err := client.Index().
		Index("todo").
		Type("_doc").
		BodyJson(todo).
		Do(reqCtx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	todo.Id = created.Id
	ctx.JSON(http.StatusCreated, todo)
}

func (s *Service) FetchTodo(ctx *gin.Context) {
	client := s.Elastic
	reqCtx := context.Background()
	todoId := ctx.Param("todo_id")
	todoRaw, err := client.Get().
		Index(todoId).
		Id(todoId).
		Do(reqCtx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !todoRaw.Found {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	var todo models.Todo
	json.Unmarshal(*todoRaw.Source, &todo)
	todo.Id = todoId
	ctx.JSON(http.StatusOK, todo)
}

func (s *Service) FetchUserTodos(ctx *gin.Context) {
	client := s.Elastic
	reqCtx := context.Background()
	query := elastic.NewTermsQuery("user_id", ctx.Param("user_id"))
	result, err := client.Search().
		Query(query).
		Index("todo").
		Type("_doc").
		Do(reqCtx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if result.Hits.TotalHits == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	todos := make([]models.Todo, 0, result.Hits.TotalHits)
	for _, hit := range result.Hits.Hits {
		var todo models.Todo
		json.Unmarshal(*hit.Source, &todo)
		todo.Id = hit.Id
		todos = append(todos, todo)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}

func (s *Service) UpdateTodo(ctx *gin.Context) {
	client := s.Elastic
	reqCtx := context.Background()
	var todo models.Todo
	ctx.BindJSON(&todo)
	res, err := client.Update().
		Index("todo").
		Type("_doc").
		Id(ctx.Param("todo_id")).
		Doc(&todo).
		FetchSource(true).
		Do(reqCtx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	json.Unmarshal(*res.GetResult.Source, &todo)
	todo.Id = res.Id
	ctx.JSON(http.StatusOK, todo)
}

func (s *Service) DeleteTodo(ctx *gin.Context) {
	client := s.Elastic
	reqCtx := context.Background()
	res, err := client.Delete().
		Index("todo").
		Id(ctx.Param("todo_id")).
		Do(reqCtx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if res.Result != "deleted" {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("invalid status: %v", res.Result))
		return
	}
	ctx.Status(http.StatusOK)
}

func (s *Service) GetUser(ctx *gin.Context) {
	client := s.Elastic
	reqCtx := context.Background()
	id := ctx.Param("user_id")
	log.Println(id)
	userRaw, err := client.Get().Index("user").Id(id).Do(reqCtx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !userRaw.Found {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	var user models.User
	user.Id = userRaw.Id
	json.Unmarshal([]byte(*userRaw.Source), &user)

	ctx.JSON(http.StatusOK, user)
}

func (s *Service) FilterUsers(ctx *gin.Context) {
	client := s.Elastic
	nickname := ctx.Query("nick")
	reqCtx := context.Background()
	termQuery := elastic.NewTermQuery("nickname", nickname)
	result, err := client.Search().
		Index("user").
		Query(termQuery).
		Do(reqCtx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if result.Hits.TotalHits == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	users := make([]models.User, 0, result.Hits.TotalHits)
	for _, hit := range result.Hits.Hits {
		var user models.User
		json.Unmarshal(*hit.Source, &user)
		users = append(users, user)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (s *Service) MakeUser(ctx *gin.Context) {
	client := s.Elastic
	reqCtx := context.Background()
	var user models.User
	ctx.BindJSON(&user)
	created, err := client.Index().
		Index("user").
		Type("_doc").
		BodyJson(user).
		Do(reqCtx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	user.Id = created.Id
	ctx.JSON(http.StatusOK, user)

}
