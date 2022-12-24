package server

import (
	"context"
	"github.com/k4zb3k/pethub/internal/repository"
	"github.com/k4zb3k/pethub/pkg/logging"
	"net/http"
)

var logger = logging.GetLogger()

const userID = "user_id"

type Auth struct {
	Repository *repository.Repository
}

type contextKey struct {
	key string
}

func (s *Server) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token := request.Header.Get("token")
		if len(token) == 0 {
			logger.Error("error")
			http.Error(writer, http.StatusText(http.StatusBadRequest), 404)
			return
		}

		userId, err := s.Services.ValidateToken(token)
		if err != nil {
			logger.Error(err)
			return
		}

		request = request.WithContext(context.WithValue(request.Context(), userID, userId))
		//logger.Info("test in WithContext", userId)

		next.ServeHTTP(writer, request)
	})

}
