package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GraysonSolus/go-rest-practice/controllers"
)

func main() {
	fmt.Println("Starting Server")

	mux := http.NewServeMux()

	mux.HandleFunc("/resize", controllers.ResizeImage)
	mux.HandleFunc("/rotate", controllers.RotateImage)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
