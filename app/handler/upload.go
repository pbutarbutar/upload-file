package handler

import (
	"brks/app/models"
	"brks/app/utils"
	"brks/domain"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var (
	validateFileType = map[string]bool{
		"image/jpeg":                true,
		"image/gif":                 true,
		"image/png":                 true,
		"application/pdf":           false,
		"image/vnd.adobe.photoshop": false,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": false,
		"text/html": false,
	}
)

type Size interface {
	Size() int64
}

func init() {
	godotenv.Load()
}

type UploadInt interface {
	HtmlUpload(w http.ResponseWriter, r *http.Request)
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

func (u *UploadHandler) HtmlUpload(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./static"))
	r.URL.Path = "/upload.html"
	fs.ServeHTTP(w, r)
}

func (u *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {

	apiResp := models.ApiResponse{
		Success: true,
		Status:  400,
		Message: "Error!",
		Data:    make(map[string]interface{}),
	}

	authToken := r.PostFormValue("auth")
	if authToken != os.Getenv("SECRET") {
		apiResp.Success = false
		apiResp.Message = "Invalid Access"
		utils.SendHTTPResponse(w, http.StatusBadRequest, apiResp)
		return
	}

	r.ParseMultipartForm(8 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		apiResp.Success = false
		apiResp.Message = fmt.Sprintf("Upload error, %s", err.Error())
		utils.SendHTTPResponse(w, http.StatusBadRequest, apiResp)
		return
	}
	defer file.Close()

	if _, ok := file.(Size); !ok {
		apiResp.Success = false
		apiResp.Message = "Fail Upload"
		utils.SendHTTPResponse(w, http.StatusBadRequest, apiResp)
		return
	}

	filetype := handler.Header.Get("Content-Type")
	if !validateFileType[filetype] {
		apiResp.Success = false
		apiResp.Message = fmt.Sprintf("invalid file type, %s", filetype)
		utils.SendHTTPResponse(w, http.StatusBadRequest, apiResp)
		return
	}

	defer file.Close()

	fileSet := "upload-*." + filepath.Ext(handler.Filename)
	tempFile, err := ioutil.TempFile("uploads", fileSet)

	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		apiResp.Success = false
		apiResp.Message = fmt.Sprintf("err: %s", err.Error())
		utils.SendHTTPResponse(w, http.StatusBadRequest, apiResp)
		return
	}

	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		apiResp.Success = false
		apiResp.Message = fmt.Sprintf("err: %s", err.Error())
		utils.SendHTTPResponse(w, http.StatusBadRequest, apiResp)
		return
	}

	tempFile.Write(fileBytes)

	var upl domain.Upload
	upl.Size = handler.Size
	upl.FileName = handler.Filename
	upl.Extension = handler.Header.Get("Content-Type")

	//Store to repository
	u.UploadEntity.Upload(r.Context(), upl)

	apiResp.Status = http.StatusOK
	apiResp.Message = "Successfully"
	utils.SendHTTPResponse(w, http.StatusOK, apiResp)

}
