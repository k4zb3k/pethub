package server

import "net/http"

func ErrResponse(w http.ResponseWriter, statusCode int, message string) {
	logger.Infoln(message)
	http.Error(w, message, statusCode)
}
