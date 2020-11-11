package server

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tunedmystic/commits.lol/app/config"
	"github.com/tunedmystic/commits.lol/app/db"
)

// Server contains all dependencies for the application.
type Server struct {
	Templates *template.Template
	DB        db.Database
}

// NewServer creates a new Server type.
func NewServer(DB db.Database) *Server {
	templateFuncs := template.FuncMap{
		"BaseURL": func() string {
			return config.App.BaseURL
		},
		"Unescape": func(html string) template.HTML {
			return template.HTML(html)
		},
	}
	s := Server{
		DB:        DB,
		Templates: template.Must(template.New("").Funcs(templateFuncs).ParseGlob("templates/*.html")),
	}
	return &s
}

// IndexHandler renders the index page.
func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Get query params and normalize.
	group := r.URL.Query().Get("group")
	fragmentParam := r.URL.Query().Get("fragment")
	if fragmentParam == "" {
		fragmentParam = "false"
	}
	fragment, _ := strconv.ParseBool(fragmentParam)

	// Get recent commits.
	commits, err := s.DB.RecentCommitsByGroup(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template.
	if fragment {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		s.Templates.ExecuteTemplate(w, "commits", commits)
		return
	}
	s.Templates.ExecuteTemplate(w, "index", commits)
}

// CacheControl ...
func CacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=300") // 5 minutes
		h.ServeHTTP(w, r)
	})
}

// Routes returns the routes for the application.
func (s *Server) Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", s.IndexHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", CacheControl(http.FileServer(http.Dir("static")))))
	return router
}
