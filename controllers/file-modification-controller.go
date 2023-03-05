package filemod

import (
	"fmt"
	"image"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
)

func resizeImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
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

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	defer file.Close()

	// TODO: move this to a security package
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	// checking the content type
	// so we don't allow files other than images
	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/jpg" {
		fmt.Println("The provided file format is not allowed. Please upload a JPEG,JPG or PNG image")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	// TODO: Move this into an image manipulation service
	rawImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Failed to decode image from file")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resized := imaging.Resize(rawImg, width, height, imaging.Lanczos)
	// err = imaging.Save(resized, fmt.Sprintf("public/single/%v", "resizedImg.png"))
	// if err != nil {
	// 	log.Fatalf("failed to save image: %v", err)
	// }

	fmt.Printf("Uploaded File: %+v\n", fileHeader.Filename)
	fmt.Printf("File Size: %+v\n", fileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", fileHeader.Header)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	imaging.Encode(w, resized, imaging.PNG)

	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
