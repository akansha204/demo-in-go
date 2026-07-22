package main

import (
	"log"

	todo "github.com/akansha204/demo-in-go"
	"github.com/akansha204/mini-rpc/rpc"
)

func main() {
	server := rpc.NewDefaultServer()

	svc := NewTodoServiceImpl()
	todo.RegisterTodoService(server, svc)

	log.Println("Todo server listening on :8080")
	if err := server.Serve(":8080"); err != nil {
		log.Fatal(err)
	}
}
