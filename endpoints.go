package lumberyard

import (
  "context"

  "github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
  PostDeploymentEndpoint   endpoint.Endpoint
  GetDeploymentsEndpoint   endpoint.Endpoint
}

// Setup Endpoints
func MakeServerEndpoints(s Store) Endpoints {
  return Endpoints {
    PostDeploymentEndpoint:  MakePostDeploymentEndpoint(s),
    GetDeploymentsEndpoint:  MakeGetDeploymentsEndpoint(s),
  }
}

// Main endpoint functions
func (e Endpoints) PostDeployment(ctx context.Context, d Deployment) error {
  request := postDeploymentRequest{Deployment: d}
  response, err := e.PostDeploymentEndpoint(ctx, request)
  if err != nil {
    return err
  }
  resp := response.(postDeploymentResponse)
  return resp.Err
}

func (e Endpoints) GetDeployments(ctx context.Context) ([]Deployment, error) {
  req := getDeploymentsRequest{}
  response, err := e.GetDeploymentsEndpoint(ctx, req)
  if err != nil {
    return nil, err
  }
  resp := response.(getDeploymentsResponse)
  return resp.Deployments, resp.Err
}

// Service Implementer functions returning endpoints
func MakePostDeploymentEndpoint(s Store) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (response interface{}, err error) {
    req := request.(postDeploymentRequest)
    e := s.PostDeployment(ctx, req.Deployment)
    return postDeploymentResponse{Err: e}, nil
  }
}

func MakeGetDeploymentsEndpoint(s Store) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (response interface{}, err error) {
    deployments, e := s.GetDeployments(ctx)
    return getDeploymentsResponse{
      Deployments: deployments,
      Err: e,
    }, nil
  }
}

// Request/Response Structs
type postDeploymentRequest struct {
  Deployment Deployment
}

type postDeploymentResponse struct {
  Err error  `json:"err,omitempty"`
}

func (r postDeploymentResponse) error() error { return r.Err }

type getDeploymentsRequest struct {}

type getDeploymentsResponse struct {
  Deployments []Deployment    `json:"deployments"`
  Err error    `json:"err,omitempty"`
}

func (r getDeploymentsResponse) error() error { return r.Err }
