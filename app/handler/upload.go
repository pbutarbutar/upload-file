package handler

import (
	"brks/domain"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
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

type UploadHandler struct {
	UploadEntity domain.UploadEntity
}

func NewUploadHandler(r *gin.RouterGroup, up domain.UploadEntity) {
	handler := &UploadHandler{
		UploadEntity: up,
	}
	r.POST("/upload", handler.UploadFile)
}

func (u *UploadHandler) UploadFile(c *gin.Context) {
	authToken := c.PostForm("auth")

	if authToken != os.Getenv("SECRET") {
		c.JSON(http.StatusForbidden, gin.H{"errors": "Auth Invalid!"})
		return
	}

	// File
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fmt.Sprintf("file err : %s", err.Error()),
		})
		return
	}

	filetype := file.Header.Get("Content-Type")

	if !validateFileType[filetype] {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fmt.Sprintf("invalid file type, %s", filetype),
		})
		return
	}

	if file.Size > maxFile { //max 8Mb
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fmt.Sprintf("file err, max size : %v Byte (8 MB)", maxFile),
		})
		return
	}

	// Set Folder untuk menyimpan filenya
	path := "public/" + file.Filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	var upl domain.Upload
	upl.Size = file.Size
	upl.FileName = file.Filename
	upl.Extension = file.Header.Get("Content-Type")

	//Store to repository
	u.UploadEntity.Upload(c.Request.Context(), upl)

	//remove file temp
	e := os.Remove(path)
	if e != nil {
		log.Error().Err(e).Msg("Cannot remove file")
	}

	c.JSON(200, gin.H{
		"data": "",
		"size": file.Size,
		"Ext":  file.Filename,
	})
}
