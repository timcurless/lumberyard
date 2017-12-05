package Pipelines

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/timcurless/lumberyard/Cassandra"
)

// Post -- handles POST request to /api/v1/pipelines to create new pipeline
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params
func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var gocqlUUID gocql.UUID

	// FormToPipeline() is included in Pipelines/processing.go
	pipeline, errs := FormToPipeline(r)

	// have we created a pipeline correctly
	var created = false

	// if we had no errors from FormToPipeline, we will
	// attempt to save our data to Cassandra
	if len(errs) == 0 {
		fmt.Println("creating a new pipeline")

		// generate a unique UUID for this pipeline
		gocqlUUID = gocql.TimeUUID()

		// write data to Cassandra
		if err := Cassandra.Session.Query(`
      INSERT INTO pipelines (id, name, description) VALUES (?, ?, ?)`,
			gocqlUUID, pipeline.Name, pipeline.Description).Exec(); err != nil {
			errs = append(errs, err.Error())
		} else {
			created = true
		}
	}

	// depending on whether we created the pipeline, return the
	// resource ID in a JSON payload, or return our errors
	if created {
		fmt.Println("pipeline_id", gocqlUUID)
		json.NewEncoder(w).Encode(NewPipelineResponse{ID: gocqlUUID})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
