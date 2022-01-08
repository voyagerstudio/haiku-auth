package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/voyagerstudio/haiku-auth/pkg/db"
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
		log.Error("empty user in getnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(userID) != 128 {
		log.Error("invalid user id in getnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	noteID := params[ParamNote]
	if noteID == "" {
		log.Error("empty note in getnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(noteID) != 128 {
		log.Error("invalid note id in getnote")
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

func (s *Server) CreateNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID := params[ParamUser]
	if userID == "" {
		log.Error("empty user in createnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(userID) != 128 {
		log.Error("invalid user id in createnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var toCreate db.Note
	err := decodeJSONBody(w, r, &toCreate)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if toCreate.ID != "" {
		msg := "Cannot create note with a specific id"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if toCreate.Text == "" {
		msg := "Cannot create note without data"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id, err := GenerateId()
	if err != nil {
		http.Error(w, "Could not store note", http.StatusInternalServerError)
		return
	}

	note, err := s.db.CreateNote(userID, id, toCreate.Text, toCreate.Order)
	if err != nil {
		// TODO: explicitly handle case of 'id already exists' error
		http.Error(w, "Could not store note", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(note)
	if err != nil {
		log.Errorf("error marshalling note %s for user %s: %v", note.ID, userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func (s *Server) UpdateNote(w http.ResponseWriter, r *http.Request) {
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

	noteID := params[ParamNote]
	if noteID == "" {
		log.Error("empty note in getnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(noteID) != 128 {
		log.Error("invalid note id in getnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var toUpdate db.Note
	err := decodeJSONBody(w, r, &toUpdate)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if toUpdate.ID != noteID {
		msg := "Payload id does not correspond to endpoint id"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if toUpdate.Text == "" {
		msg := "Cannot update a note without data"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = s.db.UpdateNote(userID, noteID, toUpdate.Text, toUpdate.Order)
	if err != nil {
		http.Error(w, "Could not store note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) DeleteNote(w http.ResponseWriter, r *http.Request) {
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

	noteID := params[ParamNote]
	if noteID == "" {
		log.Error("empty note in getnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(noteID) != 128 {
		log.Error("invalid note id in getnote")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.db.DeleteNote(userID, noteID)
	if err != nil {
		log.Error("Could not update note")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
