package controllers

import (
	"fmt"
	"image"
	"image/color"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
)

func ResizeImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Max upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	rawWidth := r.FormValue("width")
	width, err := strconv.Atoi(rawWidth)
	if err != nil {
		fmt.Println("Failed to parse width with value " + rawWidth)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	rawHeight := r.FormValue("height")
	height, err := strconv.Atoi(rawHeight)
	if err != nil {
		fmt.Println("Failed to parse height with value " + rawHeight)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if width == 0 || height == 0 {
		fmt.Println("Height or width can not be 0")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	defer file.Close()

	// TODO: Move this into an image manipulation service
	rawImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Failed to decode image from file: " + err.Error())

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resized := imaging.Resize(rawImg, width, height, imaging.Lanczos)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	imaging.Encode(w, resized, imaging.PNG)
}

func RotateImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	rawDegrees := r.FormValue("degrees")
	degrees, err := strconv.ParseFloat(rawDegrees, 32)
	if err != nil {
		fmt.Println("Failed to parse degrees with value " + rawDegrees)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	defer file.Close()

	// TODO: Move this into an image manipulation service
	rawImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Failed to decode image from file")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	rotated := imaging.Rotate(rawImg, degrees, color.Transparent)
	// err = imaging.Save(rotated, "rotatedImg.png")
	// if err != nil {
	// 	log.Fatalf("failed to save image: %v", err)
	// }

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	imaging.Encode(w, rotated, imaging.PNG)
}
