package service

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints available in this service
type Endpoints struct {
	PostProjectEndpoint      endpoint.Endpoint
	GetProjectEndpoint       endpoint.Endpoint
	PostStackEndpoint        endpoint.Endpoint
	GetProjectStacksEndpoint endpoint.Endpoint
}

// MakeServerEndpoints Create a collection of server endpoints
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostProjectEndpoint:      MakePostProjectEndpoint(s),
		GetProjectEndpoint:       MakeGetProjectEndpoint(s),
		PostStackEndpoint:        MakePostStackEndpoint(s),
		GetProjectStacksEndpoint: MakeGetProjectStacksEndpoint(s),
	}
}

// PostProject Endpoint to create a new project
func (e Endpoints) PostProject(ctx context.Context, p Project) error {
	request := postProjectRequest{Project: p}
	response, err := e.PostProjectEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postProjectResponse)
	return resp.Err
}

// GetProject endpoint to get a project
func (e Endpoints) GetProject(ctx context.Context, id string) (Project, error) {
	request := getProjectRequest{ID: id}
	response, err := e.GetProjectEndpoint(ctx, request)
	if err != nil {
		return Project{}, err
	}
	resp := response.(getProjectResponse)
	return resp.Project, resp.Err
}

// PostStack endpoint to add a new stack to a specified project
func (e Endpoints) PostStack(ctx context.Context, projectID string, s Stack) error {
	request := postStackRequest{ProjectID: projectID, Stack: s}
	response, err := e.PostStackEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postStackResponse)
	return resp.Err
}

// MakePostProjectEndpoint factory function
func MakePostProjectEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postProjectRequest)
		id, e := s.PostProject(ctx, req.Project)
		return postProjectResponse{ID: id, Err: e}, nil
	}
}

// MakeGetProjectEndpoint factory function
func MakeGetProjectEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getProjectRequest)
		p, e := s.GetProject(ctx, req.ID)
		return getProjectResponse{Project: p, Err: e}, nil
	}
}

// MakePostStackEndpoint factory function
func MakePostStackEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postStackRequest)
		id, e := s.PostStack(ctx, req.ProjectID, req.Stack)
		return postProjectResponse{ID: id, Err: e}, nil
	}
}

func MakeGetProjectStacksEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getProjectStacksRequest)
		stacks, e := s.GetProjectStacks(ctx, req.ID)
		return getProjectStacksResponse{Stacks: stacks, Err: e}, nil
	}
}

// Request/Response Types
type postProjectRequest struct {
	Project Project
}

type postProjectResponse struct {
	ID  string `json:"project_id"`
	Err error  `json:"err,omitempty"`
}

func (r postProjectResponse) error() error { return r.Err }

type getProjectRequest struct {
	ID string
}

type getProjectResponse struct {
	Project Project `json:"project,omitempty"`
	Err     error   `json:"err,omitempty"`
}

func (r getProjectResponse) error() error { return r.Err }

type postStackRequest struct {
	ProjectID string
	Stack     Stack
}

type postStackResponse struct {
	ID  string `json:"stack_id"`
	Err error  `json:"err,omitempty"`
}

func (r postStackResponse) error() error { return r.Err }

type getProjectStacksRequest struct {
	ID string
}

type getProjectStacksResponse struct {
	Stacks []Stack `json:"stacks,omitempty"`
	Err    error   `json:"err,omitempty"`
}

func (r getProjectStacksResponse) error() error { return r.Err }
