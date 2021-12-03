package main

import (
	_entity "brks/app/entity"
	_repo "brks/app/repository"
	"brks/config"

	_handler "brks/app/handler"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("unable to load env through env config")
	}
}

func main() {

	r := gin.Default()

	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Use(cors.Default())

	config.SetupModels()
	db := config.GetDBConnection()
	port := config.GetPortConnection()

	repoUpload := _repo.NewUploadRepository(db)
	entityUpload := _entity.NewUploadEntity(repoUpload)

	api := r.Group("/v1")

	_handler.NewUploadHandler(api, entityUpload)

	r.Run(port)
}
