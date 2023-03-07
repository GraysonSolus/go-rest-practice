package controller_test

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
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
		var buf bytes.Buffer
		fileWriter := multipart.NewWriter(&buf)

		err := jpeg.Encode(&buf, expectedImg, nil)
		if err != nil {
			t.Error(err)
		}

		if err = fileWriter.Close(); err != nil {
			t.Error(err)
		}

		request, err := http.NewRequest("POST", "/resize", &buf)
		if err != nil {
			t.Error(err)
		}

		request.Header.Add("Content-Type", fileWriter.FormDataContentType())

		testHeader := textproto.MIMEHeader{}
		testHeader.Add("test", "image")

		fh := &multipart.FileHeader{
			Filename: "test_img.jpeg",
			Header:   testHeader,
		}

		mfh := make([]*multipart.FileHeader, 0)
		mfh = append(mfh, fh)

		request.ParseMultipartForm(10 << 20)
		request.MultipartForm.File["image"] = mfh

		file, _, _ := request.FormFile("image") // Don't care about errors here, just need a return object for our ReadInResizeParams method

		controller := controller.NewImageModController(&fakeImageModSvc, &fakeValidator)

		fakeValidator.On("ReadInResizeParams", request).Return(file, 200, 200, nil)
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
