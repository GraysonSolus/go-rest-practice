package controller

import (
	"fmt"
	"image"
	"image/color"
	"mime/multipart"
	"net/http"

	"github.com/disintegration/imaging"
)

type imgModService interface {
	ResizeImage(file multipart.File, targetWidth, targetHeight int) (*image.NRGBA, error)
	RotateImage(file multipart.File, degrees float64, backgroundColor color.Color) (*image.NRGBA, error)
}

type imgInputValidator interface {
	ReadInResizeParams(r *http.Request) (file multipart.File, width, height int, err error)
	ReadInRotateParams(r *http.Request) (file multipart.File, degrees float64, err error)
}

type ImageModController struct {
	imageService   imgModService
	inputValidator imgInputValidator
}

func NewImageModController(imgModSvc imgModService, validator imgInputValidator) *ImageModController {
	return &ImageModController{
		imageService:   imgModSvc,
		inputValidator: validator,
	}
}

func (imc *ImageModController) HandleResizeImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	file, width, height, err := imc.inputValidator.ReadInResizeParams(r)
	if err != nil {
		fmt.Println("invalid parameters: " + err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	defer file.Close()

	resizedImg, err := imc.imageService.ResizeImage(file, width, height)
	if err != nil {
		fmt.Println("Failed to rotate image: " + err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	// TODO: handle error before we set status to OK
	imaging.Encode(w, resizedImg, imaging.PNG)
}

func (imc *ImageModController) HandleRotateImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	file, degrees, err := imc.inputValidator.ReadInRotateParams(r)
	if err != nil {
		fmt.Println("invalid parameters: " + err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	defer file.Close()

	rotated, err := imc.imageService.RotateImage(file, degrees, color.Transparent)
	if err != nil {
		fmt.Println("Failed to resize image: " + err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	// TODO: handle error before we set status to OK
	imaging.Encode(w, rotated, imaging.PNG)
}
