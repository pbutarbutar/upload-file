package handler

import (
	"brks/app/models"
	"brks/app/utils"
	"brks/domain"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var (
	validateFileType = map[string]bool{
		"image/jpeg":                true,
		"image/gif":                 true,
		"image/png":                 true,
		"application/pdf":           true,
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

	apiResp := models.ApiResponse{
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

	tempFile, err := ioutil.TempFile("static", "upload-*.png")

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

	//remove file temp
	/*filePath := "static/" + handler.Filename
	e := os.Remove(filePath)
	if e != nil {
		log.Error().Err(e).Msg("Cannot remove file")
	}*/

	utils.SendHTTPResponse(w, http.StatusBadRequest, apiResp)

}
