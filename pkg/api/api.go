package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	rand "crypto/rand"

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

	r.HandleFunc(fmt.Sprintf("/user/{%s}/notes", ParamUser), s.GetNoteList).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/user/{%s}/notes", ParamUser), s.CreateNote).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/user/{%s}/note/{%s}", ParamUser, ParamNote), s.GetNote).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("/user/{%s}/note/{%s}", ParamUser, ParamNote), s.DeleteNote).Methods(http.MethodDelete)
	r.HandleFunc(fmt.Sprintf("/user/{%s}/note/{%s}", ParamUser, ParamNote), s.UpdateNote).Methods(http.MethodPut)

	// r.HandleFunc("/user", s.CreateUser).Methods(http.MethodPost)
	// r.HandleFunc(fmt.Sprintf("/user/{%s}", ParamUser), s.GetUser).Methods(http.MethodGet)

	s.srv.Handler = r

	return s
}

const ID_SIZE = 128

func GenerateId() (string, error) {
	raw := make([]byte, ID_SIZE)
	_, err := rand.Read(raw)

	if err != nil {
		return "", err
	}

	return string(raw), nil
}

// ListenAndServe begins listening on the designated port and serving requests
func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// Code below stolen from https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

// ParseValueAndParams parses a comma separated list of values with optional
// semicolon separated name-value pairs. Content-Type and Content-Disposition
// headers are in this format.
func ParseValueAndParams(header http.Header, key string) (value string, params map[string]string) {
	params = make(map[string]string)
	s := header.Get(key)
	value, s = expectTokenSlash(s)
	if value == "" {
		return
	}
	value = strings.ToLower(value)
	s = skipSpace(s)
	for strings.HasPrefix(s, ";") {
		var pkey string
		pkey, s = expectToken(skipSpace(s[1:]))
		if pkey == "" {
			return
		}
		if !strings.HasPrefix(s, "=") {
			return
		}
		var pvalue string
		pvalue, s = expectTokenOrQuoted(s[1:])
		if pvalue == "" {
			return
		}
		pkey = strings.ToLower(pkey)
		params[pkey] = pvalue
		s = skipSpace(s)
	}
	return
}

// Octet types from RFC 2616.
type octetType byte

var octetTypes [256]octetType

const (
	isToken octetType = 1 << iota
	isSpace
)

func skipSpace(s string) (rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isSpace == 0 {
			break
		}
	}
	return s[i:]
}

func expectToken(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		if octetTypes[s[i]]&isToken == 0 {
			break
		}
	}
	return s[:i], s[i:]
}

func expectTokenSlash(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		b := s[i]
		if (octetTypes[b]&isToken == 0) && b != '/' {
			break
		}
	}
	return s[:i], s[i:]
}

func expectTokenOrQuoted(s string) (value string, rest string) {
	if !strings.HasPrefix(s, "\"") {
		return expectToken(s)
	}
	s = s[1:]
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			return s[:i], s[i+1:]
		case '\\':
			p := make([]byte, len(s)-1)
			j := copy(p, s[:i])
			escape := true
			for i = i + 1; i < len(s); i++ {
				b := s[i]
				switch {
				case escape:
					escape = false
					p[j] = b
					j++
				case b == '\\':
					escape = true
				case b == '"':
					return string(p[:j]), s[i+1:]
				default:
					p[j] = b
					j++
				}
			}
			return "", ""
		}
	}
	return "", ""
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}
