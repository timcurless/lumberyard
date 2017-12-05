package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/timcurless/lumberyard/Cassandra"
	"github.com/timcurless/lumberyard/Pipelines"
	"github.com/timcurless/lumberyard/Stages"
	"github.com/timcurless/lumberyard/Projects"
)

type heartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func main() {
	CassandraSession := Cassandra.Session
	defer CassandraSession.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)
	router.HandleFunc("/api/v1/pipelines", Pipelines.Get).Methods("GET")
	router.HandleFunc("/api/v1/pipelines", Pipelines.Post).Methods("POST")
	router.HandleFunc("/api/v1/pipelines/{pipeline_uuid}", Pipelines.GetOne).Methods("GET")
	router.HandleFunc("/api/v1/pipelines/{pipeline_uuid}/stages", Stages.Post).Methods("POST")
	router.HandleFunc("/api/v1/projects", Projects.Post).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", router))
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}
