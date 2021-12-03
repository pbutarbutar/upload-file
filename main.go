package main

import (
	_entity "brks/app/entity"
	_router "brks/app/http"
	_repo "brks/app/repository"
	"brks/config"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_handler "brks/app/handler"

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
	config.SetupModels()
	db := config.GetDBConnection()
	port := config.GetPortConnection()

	repoUpload := _repo.NewUploadRepository(db)
	entityUpload := _entity.NewUploadEntity(repoUpload)
	handlerload := _handler.NewUploadHandler(entityUpload)
	routes := _router.GetRouters(handlerload)
	srv := &http.Server{
		Addr:    port,
		Handler: routes,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("fail to start service")
		}
	}()
	log.Info().Str("address", port).Msg("service started")
	<-done
	log.Info().Msg("service is going to stop")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Msg("service shutdown failed")
	}
	log.Info().Msg("service exited properly")
}
