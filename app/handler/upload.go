package handler

import (
	"brks/app/utils"
	"brks/domain"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

var maxFile int64 = 8388608
var validateFileType = map[string]bool{
	"image/jpeg":      true,
	"image/gif":       true,
	"image/png":       true,
	"application/pdf": false,
}

type UploadInt interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}

type UploadHandler struct {
	UploadEntity domain.UploadEntity
}

func NewUploadHandler(up domain.UploadEntity) UploadInt {
	return &UploadHandler{
		UploadEntity: up,
	}

}

func (u *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {

	apiResp := domain.ApiResponse{
		Success: true,
		Status:  400,
		Message: "Successfully",
		Data:    make(map[string]interface{}),
	}

	authToken := r.PostFormValue("auth")
	if authToken != os.Getenv("SECRET") {
		apiResp.Success = false
		apiResp.Message = "Invalid Access"
		utils.SendHTTPResponse(w, http.StatusBadRequest, apiResp)
		return
	}

	utils.SendHTTPResponse(w, http.StatusBadRequest, apiResp)

}
