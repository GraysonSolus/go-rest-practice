package main

import (
	"fmt"
	"log"
	"net/http"

	controllers "github.com/GraysonSolus/go-rest-practice/internal/image-modification-service/controller"
	imgmod "github.com/GraysonSolus/go-rest-practice/internal/image-modification-service/svc"
	"github.com/GraysonSolus/go-rest-practice/internal/image-modification-service/validator"
)

func main() {
	fmt.Println("Starting Server")

	imgModSvc := &imgmod.ImageModificationService{}
	inputValidator := validator.NewImageInputValidator(1, 1000, 1, 1000)

	controller := controllers.NewImageModController(imgModSvc, inputValidator)

	mux := http.NewServeMux()

	mux.HandleFunc("/resize", controller.HandleResizeImage)
	mux.HandleFunc("/rotate", controller.HandleRotateImage)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
