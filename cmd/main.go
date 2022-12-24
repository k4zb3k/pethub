package main

import (
	"github.com/gorilla/mux"
	"github.com/k4zb3k/pethub/config"
	"github.com/k4zb3k/pethub/internal/db"
	"github.com/k4zb3k/pethub/internal/repository"
	"github.com/k4zb3k/pethub/internal/server"
	"github.com/k4zb3k/pethub/internal/services"
	"github.com/k4zb3k/pethub/pkg/logging"
	"net"
	"net/http"
)

var logger = logging.GetLogger()

func main() {
	err := execute()
	if err != nil {
		logger.Error(err)
	}
}

func execute() error {
	router := mux.NewRouter()

	connection, err := db.GetDbConnection()
	if err != nil {
		logger.Error("error connecting to DB", err)
	}
	newRepository := repository.NewRepository(connection)

	service := services.NewServices(newRepository)

	newServer := server.NewServer(router, service)

	newServer.Init()

	getConfig, err := config.GetConfig()
	if err != nil {
		logger.Error("GetConfig is crashed", err)
	}
	address := net.JoinHostPort(getConfig.Host, getConfig.Port)
	srv := http.Server{
		Addr:    address,
		Handler: router,
	}
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("error in ListenAndServe", err)
	}

	return nil
}
