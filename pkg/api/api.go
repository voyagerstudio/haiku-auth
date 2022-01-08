package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/voyagerstudio/haiku-auth/pkg/db"
)

const (
	ParamUser = "user"
	ParamNote = "note"
)

// Server is a wrapper type for the general HTTP server
// We'll be adding things in here like references to a database
type Server struct {
	srv *http.Server
	db  *db.Conn
}

// NewServer instantiates a new HTTP REST server
func NewServer(host string, port int, db *db.Conn) *Server {
	s := &Server{
		srv: &http.Server{
			Addr: fmt.Sprintf("%s:%d", host, port),
			// Default timeouts are unlim, which is bad
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		db: db,
	}

	// We could use the stdlib muxer, but gorilla is incredibly nice,
	// lightweight, fulfills the standard interfaces, and comes with some
	// nice additional features
	r := mux.NewRouter()
	r.HandleFunc("/ping", s.PingHandler)

	r.HandleFunc(fmt.Sprintf("/note/{%s}", ParamNote), s.GetNote).Methods(http.MethodGet)

	r.HandleFunc(fmt.Sprintf("/user/{%s}/notes", ParamUser), s.GetNoteList).Methods(http.MethodGet)

	s.srv.Handler = r

	return s
}

// ListenAndServe begins listening on the designated port and serving requests
func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}
