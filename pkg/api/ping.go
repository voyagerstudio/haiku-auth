package api

import "net/http"

// PingHandler serves an HTTP 200
func (s *Server) PingHandler(w http.ResponseWriter, r *http.Request) {
	// 200 is the default if we return without writing a header,
	// but it's nice to be verbose. This *could* just be an empty
	// function and it'd act the same.
	w.WriteHeader(http.StatusOK)
}
