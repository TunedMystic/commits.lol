package server

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/tunedmystic/commits.lol/app/db"
)

// Server contains all dependencies for the application.
type Server struct {
	Templates *template.Template
	Router    *mux.Router
	DB        db.Database
}

// NewServer creates a new Server type.
func NewServer(DB db.Database) *Server {
	s := Server{}
	s.DB = DB
	s.Templates = template.Must(template.New("").ParseGlob("templates/*.html"))
	return &s
}

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	commits, err := s.DB.RecentCommits()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s.Templates.ExecuteTemplate(w, "index", commits)
}

func (s *Server) Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", s.IndexHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return router
}

type TemplateContext struct {
	//
}
