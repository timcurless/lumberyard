package Pipelines

import (
	"github.com/gocql/gocql"
)

// Pipeline Struct to hold info about a pipeline
type Pipeline struct {
	ID          gocql.UUID `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}

// GetPipelineResponse - Returns a pipeline
type GetPipelineResponse struct {
	Pipeline Pipeline `json:"pipeline"`
}

// AllPipelinesResponse - Returns a list of all pipelines
type AllPipelinesResponse struct {
	Pipelines []Pipeline `json:"pipelines"`
}

// NewPipelineResponse - returns ID of newly created pipeline
type NewPipelineResponse struct {
	ID gocql.UUID `json:"id"`
}

// ErrorResponse - returns error if applicable
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
