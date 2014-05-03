package main

import (
	"fmt"
	"os"

	"github.com/walkr/kvlite3"
)

type Todo struct {
	Key   string
	Title string
}

func main() {

	client, _ := kvlite3.NewClient("./todos.db")
	bucket, _ := client.NewBucket("todos")

	todo := Todo{kvlite3.NewKey(), "Take out trash"}
	bucket.Put(todo.Key, todo) // Insert
	fmt.Println("Todo saved")

	todo2 := Todo{}
	bucket.Get(todo.Key, &todo2) // Fetch
	fmt.Printf("Todo is: %v\n", todo2)

	todos := []*Todo{}
	bucket.GetAll(&todos) // Fetch all
	fmt.Printf("Todos are: %v\n", todos)

	os.Remove("./todos.db")
}
