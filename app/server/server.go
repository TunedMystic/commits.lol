package server

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/getsentry/sentry-go"
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
		"Goatcounter": func() string {
			return config.App.GoatcounterUser
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
	// Should we render the recent commits as an html fragment?
	fragmentParam := r.URL.Query().Get("fragment")
	if fragmentParam == "" {
		fragmentParam = "false"
	}
	fragment, _ := strconv.ParseBool(fragmentParam)

	// Get recent commits.
	group := r.URL.Query().Get("group")
	commits, err := s.DB.RecentCommitsByGroup(group)
	if err != nil {
		sentry.CaptureException(err)
		http.Error(w, "oopsie, something went horribly wrong", http.StatusInternalServerError)
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

// Routes returns the routes for the application.
func (s *Server) Routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", s.IndexHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", CacheControl(http.FileServer(http.Dir("static")))))
	return Logging(router)
}
