package Stages

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
  "github.com/gorilla/mux"
	"github.com/timcurless/lumberyard/Cassandra"
)

// Post -- handles POST request to /api/v1/pipelines/{pipeline_uuid}/ to create new stage
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params
func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var gocqlUUID gocql.UUID
  var created = false

  vars := mux.Vars(r)
	pipelineID := vars["pipeline_uuid"]

  pipelineUUID, err := gocql.ParseUUID(pipelineID)
  if err != nil {
		errs = append(errs, err.Error())
	} else {

  	stage, errs := FormToStage(r)

    stage.PipelineID = pipelineUUID


  	// if we had no errors from FormToStage, we will
  	// attempt to save our data to Cassandra
  	if len(errs) == 0 {
  		fmt.Println("creating a new stage for pipeline " + (pipelineUUID).String())

  		// generate a unique UUID for this pipeline
  		gocqlUUID = gocql.TimeUUID()

  		// write data to Cassandra
  		if err := Cassandra.Session.Query(`
        INSERT INTO stages (id, pipeline_id, name, description, type, version, payload) VALUES (?, ?, ?, ?, ?, ?, ?)`,
  			gocqlUUID, stage.PipelineID, stage.Name, stage.Description, stage.Type, stage.Version, stage.Payload).Exec(); err != nil {
  			errs = append(errs, err.Error())
  		} else {
  			created = true
  		}
    }
	}

	// depending on whether we created the stage, return the
	// resource ID in a JSON payload, or return our errors
	if created {
		fmt.Println("stage_id", gocqlUUID)
		json.NewEncoder(w).Encode(NewStageResponse{ID: gocqlUUID})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
