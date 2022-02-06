package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// GetUser returns all note information for a given note ID
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params[ParamUser]
	if userID == "" {
		log.Error("empty user in getnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(userID) != 128 {
		log.Error("invalid user id in getnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUser(userID)
	if err != nil {
		log.Errorf("error getting user %s: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Errorf("error marshalling user %s: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	id, err := GenerateId()
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	user, err := s.db.CreateUser(id)
	if err != nil {
		// TODO: explicitly handle case of 'id already exists' error
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Errorf("error marshalling user %s: %v", user.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}
