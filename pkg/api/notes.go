package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// GetNoteList returns a list of note IDs for a given user ID
func (s *Server) GetNoteList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params[ParamUser]
	if user == "" {
		log.Error("empty user in getnotelist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notes, err := s.db.GetNoteList(user)
	if err != nil {
		log.Errorf("error getting note list for user %s: %v", user, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(notes)
	if err != nil {
		log.Errorf("error marshalling note list for user %s: %v", user, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

// GetNote returns all note information for a given note ID
func (s *Server) GetNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params[ParamUser]
	if userID == "" {
		log.Error("empty user in getnotelist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	noteID := params[ParamNote]
	if noteID == "" {
		log.Error("empty note in getnotelist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	note, err := s.db.GetNote(userID, noteID)
	if err != nil {
		log.Errorf("error getting note %s for user %s: %v", noteID, userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(note)
	if err != nil {
		log.Errorf("error marshalling note %s for user %s: %v", noteID, userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
