package imgmod

import (
	"fmt"
	"image"
	"image/color"
	"mime/multipart"

	"github.com/disintegration/imaging"
)

type ImageModificationService struct{}

func (ims *ImageModificationService) ResizeImage(file multipart.File, targetWidth, targetHeight int) (*image.NRGBA, error) {
	rawImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Failed to decode image from file: " + err.Error())

		return nil, err
	}

	resized := imaging.Resize(rawImg, targetWidth, targetHeight, imaging.Lanczos)

	return resized, nil
}

func (ims *ImageModificationService) RotateImage(file multipart.File, degrees float64, backgroundColor color.Color) (*image.NRGBA, error) {
	rawImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Failed to decode image from file")

		return nil, err
	}

	rotated := imaging.Rotate(rawImg, degrees, backgroundColor)

	return rotated, nil
}
