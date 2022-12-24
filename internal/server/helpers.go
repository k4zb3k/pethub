package server

import (
	"net/http"
)

func BadRequest(w http.ResponseWriter, err error) {
	logger.Infof("%+v\n", err)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter, err error) {
	logger.Infof("%v\n", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
