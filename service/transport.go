package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting error for route programming issue
	ErrBadRouting = errors.New("incosistent mapping between route and handler")
)

// MakeHTTPHandler Creates an HTTP Handler of all endpoints in the service
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {

	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST /projects/              Adds a new Project
	// GET /projects/:id            Gets a Project
	// POST /projects/:id/stacks/		Add a new Stack to a project

	r.Methods("POST").Path("/projects/").Handler(httptransport.NewServer(
		e.PostProjectEndpoint,
		decodePostProjectRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/projects/{id}").Handler(httptransport.NewServer(
		e.GetProjectEndpoint,
		decodeGetProjectRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/projects/{id}/stacks").Handler(httptransport.NewServer(
		e.PostStackEndpoint,
		decodePostStackRequest,
		encodeResponse,
		options...,
	))

	return r
}

// Decode Request Functions
func decodePostProjectRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postProjectRequest

	randUUID, e := gocql.RandomUUID()
	if e != nil {
		return "", e
	}
	randUUIDstr := randUUID.String()

	if e := json.NewDecoder(r.Body).Decode(&req.Project); e != nil {
		return nil, e
	}
	req.Project.CreatedTs = time.Now().String()
	req.Project.UpdateTs = time.Now().String()
	req.Project.ID = randUUIDstr
	return req, nil
}

func decodeGetProjectRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getProjectRequest{ID: id}, nil
}

func decodePostStackRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var stack Stack
	if err := json.NewDecoder(r.Body).Decode(&stack); err != nil {
		return nil, err
	}
	return postStackRequest{
		ProjectID: id,
		Stack:     stack,
	}, nil
}

// Encode Response Functions
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrAlreadyExists:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
