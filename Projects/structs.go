package Projects

import (
  "time"
  
  "github.com/gocql/gocql"
)

type Project struct {
  ID            gocql.UUID `json:"id"`
  Name          string     `json:"name"`
  Email         string     `json:"email"`
  UpdateTS      time.Time  `json:"updated_ts"`
  CreatedTS     time.Time  `json:"created_ts"`
}

// NewProjectResponse - returns ID of newly created project
type NewProjectResponse struct {
	ID gocql.UUID `json:"id"`
}

// GetProjectResponse - Returns a Project
type GetProjectResponse struct {
	Project Project `json:"project"`
}

// ErrorResponse - returns error if applicable
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
