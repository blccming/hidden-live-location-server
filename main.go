package main

import (
	"fmt"
)

func main() {
	fmt.Println("test")
	r := initEndpoints()
	r.Run("localhost:8080")
}
