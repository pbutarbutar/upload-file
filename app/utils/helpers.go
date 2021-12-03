package utils

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func SendHTTPResponse(w http.ResponseWriter, statusCode int, model interface{}) {
	json, err := json.Marshal(model)
	if err != nil {
		log.Err(err).Msg("something went wrong")
	}
	w.WriteHeader(statusCode)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
