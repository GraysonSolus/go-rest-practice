package validator

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
)

type ImageInputValidator struct {
	minHeight int
	maxHeight int
	minWidth  int
	maxWidth  int
}

func NewImageInputValidator(minHeight, maxHeight, minWidth, maxWidth int) *ImageInputValidator {
	return &ImageInputValidator{
		minHeight: minHeight,
		maxHeight: maxHeight,
		minWidth:  minWidth,
		maxWidth:  maxWidth,
	}
}

func (iv *ImageInputValidator) ReadInResizeParams(r *http.Request) (file multipart.File, width, height int, err error) {
	r.ParseMultipartForm(10 << 20)

	rawWidth := r.FormValue("width")
	width, err = strconv.Atoi(rawWidth)
	if err != nil {
		return nil, 0, 0, errors.New("error retrieving width data from request: " + err.Error())
	}

	rawHeight := r.FormValue("height")
	height, err = strconv.Atoi(rawHeight)
	if err != nil {
		return nil, 0, 0, errors.New("error retrieving height data from request: " + err.Error())
	}

	if paramsAreValid := iv.ValidateResizeParams(width, height); !paramsAreValid {
		// TODO: pass back more usable info to end user about valid values
		return nil, 0, 0, errors.New("width and height are outside allowed bounds")
	}

	file, _, err = r.FormFile("image")
	if err != nil {
		return nil, 0, 0, errors.New("error retrieving file data from request: " + err.Error())
	}

	return file, width, height, err
}

func (iv *ImageInputValidator) ReadInRotateParams(r *http.Request) (file multipart.File, degrees float64, err error) {
	r.ParseMultipartForm(10 << 20)

	rawDegrees := r.FormValue("degrees")
	degrees, err = strconv.ParseFloat(rawDegrees, 32)
	if err != nil {
		return nil, 0, errors.New("error retrieving degrees of rotation data from request: " + err.Error())
	}

	if paramsAreValid := iv.ValidateRotateParams(degrees); !paramsAreValid {
		// TODO: pass back more usable info to end user about valid values
		return nil, 0, errors.New("expected degrees to be between 0 and 360 but recieved: " + strconv.Itoa(int(degrees))) // Quick and messy degrees to string transofrmation - consider using fmt
	}

	file, _, err = r.FormFile("image")
	if err != nil {
		return nil, 0, errors.New("error retrieving file data from request: " + err.Error())
	}

	return file, degrees, err
}

func (iv *ImageInputValidator) ValidateResizeParams(width, height int) bool {
	if width > iv.maxWidth || width < iv.minWidth {
		return false
	}

	if height > iv.maxHeight || height < iv.minHeight {
		return false
	}

	return true
}

func (iv *ImageInputValidator) ValidateRotateParams(degrees float64) bool {
	if degrees > 360 || degrees <= 0 {
		return false
	}

	return true
}
