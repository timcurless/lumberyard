package Pipelines

import (
	"encoding/json"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/timcurless/lumberyard/Cassandra"
)

// Get -- handles GET request to /api/v1/pipelines to fetch all pipelines
// params:
// w - response writer for building jSON payload Response
// r - request reader to fetch form data or url params (unused here)
func Get(w http.ResponseWriter, r *http.Request) {
	var pipelineList []Pipeline
	m := map[string]interface{}{}

	query := "SELECT id,name,description FROM pipelines"
	iterable := Cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		pipelineList = append(pipelineList, Pipeline{
			ID:          m["id"].(gocql.UUID),
			Name:        m["name"].(string),
			Description: m["description"].(string),
		})
		m = map[string]interface{}{}
	}

	json.NewEncoder(w).Encode(AllPipelinesResponse{Pipelines: pipelineList})
}

// GetOne -- handles GET request to /api/v1/pipelines/{pipeline_uuid} to fetch one pipeline
// params:
// w - response writer for building jSON payload Response
// r - request reader to fetch form data or url params
func GetOne(w http.ResponseWriter, r *http.Request) {
	var pipeline Pipeline
	var errs []string
	var found = false

	vars := mux.Vars(r)
	id := vars["pipeline_uuid"]

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		m := map[string]interface{}{}
		query := "SELECT id,name,description FROM pipelines WHERE id=? LIMIT 1"
		iterable := Cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
		for iterable.MapScan(m) {
			found = true
			pipeline = Pipeline{
				ID:          m["id"].(gocql.UUID),
				Name:        m["name"].(string),
				Description: m["description"].(string),
			}
		}
		if !found {
			errs = append(errs, "Pipeline not found")
		}
	}

	if found {
		json.NewEncoder(w).Encode(GetPipelineResponse{Pipeline: pipeline})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
