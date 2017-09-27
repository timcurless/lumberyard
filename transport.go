package lumberyard

import (
  "context"
  "encoding/json"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/go-kit/kit/log"
  httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler (s Store, logger log.Logger) http.Handler {
  r := mux.NewRouter()
  e := MakeServerEndpoints(s)
  options := []httptransport.ServerOption{
    httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

  // POST     /v1/deployments/             adds another deployment to Store
  // GET      /v1/deployments/             returns all deployments from Store

  r.Methods("POST").Path("/v1/deployments/").Handler(httptransport.NewServer(
    e.PostDeploymentEndpoint,
    decodePostDeploymentRequest,
    encodeResponse,
    options...,
  ))
  r.Methods("GET").Path("/v1/deployments/").Handler(httptransport.NewServer(
    e.GetDeploymentsEndpoint,
    decodeGetDeploymentsRequest,
    encodeResponse,
    options...,
  ))
  return r
}

func decodePostDeploymentRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
  var req postDeploymentRequest
  if e := json.NewDecoder(r.Body).Decode(&req.Deployment); e != nil {
    return nil, e
  }
  return req, nil
}

func decodeGetDeploymentsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
  return getDeploymentsRequest{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
  if e, ok := response.(errorer); ok && e.error() != nil {
    // business logic error
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
    panic("errorEncode with nil error")
  }
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(codeFrom(err))
  json.NewEncoder(w).Encode(map[string]interface{}{
    "error": err.Error(),
  })
}

func codeFrom(err error) int {
    return http.StatusInternalServerError
}
