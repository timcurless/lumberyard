package Pipelines

import (
	"encoding/json"
	"net/http"

	"github.com/gocql/gocql"
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
