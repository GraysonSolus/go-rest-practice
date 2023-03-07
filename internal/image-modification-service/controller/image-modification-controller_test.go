package controller_test

import (
	"image"
	"image/color"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GraysonSolus/go-rest-practice/internal/image-modification-service/controller"
	"github.com/stretchr/testify/mock"
)

type mockValidator struct {
	mock.Mock
}

func (mv *mockValidator) ReadInResizeParams(r *http.Request) (file multipart.File, width, height int, err error) {
	args := mv.Called(r)
	return args.Get(0).(multipart.File), args.Int(1), args.Int(2), args.Error(3)
}

func (mv *mockValidator) ReadInRotateParams(r *http.Request) (file multipart.File, degrees float64, err error) {
	args := mv.Called(r)
	return args.Get(0).(multipart.File), args.Get(1).(float64), args.Error(2)
}

type mockImageModService struct {
	mock.Mock
}

func (mims *mockImageModService) ResizeImage(file multipart.File, targetWidth, targetHeight int) (*image.NRGBA, error) {
	args := mims.Called(file, targetWidth, targetHeight)
	return args.Get(0).(*image.NRGBA), args.Error(1)
}

func (mims *mockImageModService) RotateImage(file multipart.File, degrees float64, backgroundColor color.Color) (*image.NRGBA, error) {
	args := mims.Called(file, degrees, backgroundColor)
	return args.Get(0).(*image.NRGBA), args.Error(1)
}

var (
	fakeImageModSvc mockImageModService
	fakeValidator   mockValidator
)

func setupTestCase(t *testing.T) {
	fakeImageModSvc = mockImageModService{}
	fakeValidator = mockValidator{}
}

func TestImageModController_HandleResizeImageOperation(t *testing.T) {
	t.Run("Successfully returns 200 if no error returned from service layer", func(t *testing.T) {
		setupTestCase(t)

		expectedImg := &image.NRGBA{}

		// pr, pw := io.Pipe()

		// fileWriter := multipart.NewWriter(pw)

		// defer fileWriter.Close()

		// filePart, err := fileWriter.CreateFormFile("image", "test_img.jpeg")
		// if err != nil {
		// 	t.Error(err)
		// }

		// err = jpeg.Encode(filePart, expectedImg, nil)
		// if err != nil {
		// 	t.Error(err)
		// }

		// fileWriter.Close()

		request, err := http.NewRequest("POST", "/resize", nil)
		if err != nil {
			t.Error(err)
		}

		request.MultipartForm.File["image"] = nil
		// request.Header.Add("Content-Type", fileWriter.FormDataContentType())

		controller := controller.NewImageModController(&fakeImageModSvc, &fakeValidator)

		fakeValidator.On("ReadInResizeParams", request).Return(nil, 200, 200, nil)
		fakeImageModSvc.On("ResizeImage", mock.Anything, mock.Anything, mock.Anything).Return(expectedImg, nil)

		resRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.HandleResizeImage)

		handler.ServeHTTP(resRecorder, request)

		if respStatus := resRecorder.Code; respStatus != http.StatusOK {
			t.Errorf("Expected %v but got %v", http.StatusOK, respStatus)
		}

		fakeValidator.AssertNumberOfCalls(t, "ReadInResizeParams", 1)
		fakeImageModSvc.AssertNumberOfCalls(t, "ResizeImage", 1)
	})
}
