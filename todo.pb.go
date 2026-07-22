package todo

import (
	"github.com/akansha204/mini-rpc/rpc"
)

type TodoRequest struct {
	Title string
	Description string
	Id int32
}

type TodoResponse struct {
	Id int32
	Title string
	Description string
	Completed bool
}

type TodoService interface {
	CreateTodo(req TodoRequest) (TodoResponse, error)
	GetTodo(req TodoRequest) (TodoResponse, error)
}

type TodoServiceClient struct {
}

func NewTodoServiceClient() *TodoServiceClient {
	return &TodoServiceClient{}
}

func (c *TodoServiceClient) CreateTodo(req TodoRequest) (TodoResponse, error) {
	panic("not implemented")
}

func (c *TodoServiceClient) GetTodo(req TodoRequest) (TodoResponse, error) {
	panic("not implemented")
}

func RegisterTodoService(server *rpc.Server, service TodoService) {
	server.Register(
		"TodoService/CreateTodo",
		func(payload []byte) ([]byte, error) {

			var req TodoRequest

			if err := server.Decode(payload, &req); err != nil {
				return nil, err
			}

			resp, err := service.CreateTodo(req)
			if err != nil {
				return nil, err
			}

			return server.Encode(resp)
		},
	)

	server.Register(
		"TodoService/GetTodo",
		func(payload []byte) ([]byte, error) {

			var req TodoRequest

			if err := server.Decode(payload, &req); err != nil {
				return nil, err
			}

			resp, err := service.GetTodo(req)
			if err != nil {
				return nil, err
			}

			return server.Encode(resp)
		},
	)

}

