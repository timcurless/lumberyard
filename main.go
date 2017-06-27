package main

/*
ccm create -v 3.1.1 lumberyard -n 1 -s
echo "use lumberyard; CREATE TABLE pipelines ( id UUID, name text, description text, PRIMARY KEY (id));" | cqlsh --version 3.4.2
echo "use lumberyard; CREATE TABLE stages ( id UUID, pipeline_id UUID, name text, description text, type text, version int, payload text, PRIMARY KEY (id));" | cqlsh --version 3.4.2
echo "use lumberyard; create index on stages (pipeline_id);" | cqlsh --version 3.4.2
*/

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/timcurless/lumberyard/Cassandra"
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

	log.Fatal(http.ListenAndServe(":8080", router))
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}
