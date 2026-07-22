package main

import (
	"fmt"
	"log"

	todo "github.com/akansha204/demo-in-go"
	"github.com/akansha204/mini-rpc/rpc"
)

func main() {
	client, err := rpc.NewDefaultClient(":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create a todo
	createReq := todo.TodoRequest{
		Title:       "Learn mini-gRPC",
		Description: "Build a demo using mini-protoc and mini-rpc",
	}
	var createResp todo.TodoResponse
	if err := client.Call("TodoService/CreateTodo", createReq, &createResp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[client] created: id=%d title=%q completed=%v\n",
		createResp.Id, createResp.Title, createResp.Completed)

	// Create another todo
	createReq2 := todo.TodoRequest{
		Title:       "Write Notion doc",
		Description: "Document the project for recruiters",
	}
	var createResp2 todo.TodoResponse
	if err := client.Call("TodoService/CreateTodo", createReq2, &createResp2); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[client] created: id=%d title=%q completed=%v\n",
		createResp2.Id, createResp2.Title, createResp2.Completed)

	// Get the first todo
	getReq := todo.TodoRequest{Id: createResp.Id}
	var getResp todo.TodoResponse
	if err := client.Call("TodoService/GetTodo", getReq, &getResp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[client] fetched: id=%d title=%q description=%q\n",
		getResp.Id, getResp.Title, getResp.Description)
}
