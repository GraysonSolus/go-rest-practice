package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}
