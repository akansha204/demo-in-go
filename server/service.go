package main

import (
	"fmt"
	"sync"

	todo "github.com/akansha204/demo-in-go"
)

type TodoServiceImpl struct {
	mu     sync.Mutex
	todos  map[int32]todo.TodoResponse
	nextID int32
}

func NewTodoServiceImpl() *TodoServiceImpl {
	return &TodoServiceImpl{
		todos:  make(map[int32]todo.TodoResponse),
		nextID: 1,
	}
}

func (s *TodoServiceImpl) CreateTodo(req todo.TodoRequest) (todo.TodoResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := todo.TodoResponse{
		Id:          s.nextID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}
	s.todos[t.Id] = t
	s.nextID++

	fmt.Printf("[server] created todo %d: %s\n", t.Id, t.Title)
	return t, nil
}

func (s *TodoServiceImpl) GetTodo(req todo.TodoRequest) (todo.TodoResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.todos[req.Id]
	if !ok {
		return todo.TodoResponse{}, fmt.Errorf("todo %d not found", req.Id)
	}

	fmt.Printf("[server] fetched todo %d: %s\n", t.Id, t.Title)
	return t, nil
}
