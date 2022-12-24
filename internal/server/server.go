package server

import (
	"github.com/gorilla/mux"
	"github.com/k4zb3k/pethub/internal/services"
	"net/http"
)

type Server struct {
	Mux      *mux.Router
	Services *services.Services
}

func NewServer(mux *mux.Router, services *services.Services) *Server {
	return &Server{
		Mux:      mux,
		Services: services,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

//============================================================

func (s *Server) Init() {
	logAuth := mux.MiddlewareFunc(s.ValidateToken)
	generalRout := s.Mux.PathPrefix("/api/v1").Subrouter()

	authRoute := generalRout.PathPrefix("/auth").Subrouter()
	authRoute.HandleFunc("/registration", s.Registration).Methods("POST")
	authRoute.HandleFunc("/login", s.Login).Methods("POST")

	test := generalRout.PathPrefix("/ads").Subrouter()
	test.Use(logAuth)
	test.HandleFunc("/ad", s.AddNewAd).Methods("POST")
	test.HandleFunc("/ad", s.EditAd).Methods("PUT")
	test.HandleFunc("/ad", s.DeleteAd).Methods("DELETE")
	test.HandleFunc("/my", s.GetMyAds).Methods("GET")

	s.Mux.HandleFunc("/paginate", s.Paginate).Methods("GET")
	s.Mux.HandleFunc("/search", s.Search).Methods("GET")
	s.Mux.HandleFunc("/filter", s.Filter).Methods("GET")

}
