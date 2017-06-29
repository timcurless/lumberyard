package Stages

import (
	"github.com/gocql/gocql"
)

// Stage Struct to hold info about a stage
type Stage struct {
	ID          gocql.UUID `json:"id"`
	PipelineID  gocql.UUID `json:"pipeline_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Version     int        `json:"version"`
	Payload     string     `json:"payload"`
}

// GetStageResponse - Returns a stage
type GetStageResponse struct {
	Stage Stage `json:"stage"`
}

// AllStagesForPipelineResponse - Returns a list of all stages for a pipeline
type AllStagesForPipelineResponse struct {
	Stages []Stage `json:"stages"`
}

// NewStageResponse - returns ID of newly created stage
type NewStageResponse struct {
	ID gocql.UUID `json:"id"`
}

// ErrorResponse - returns error if applicable
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
