package Projects

import (
  "encoding/json"
  "fmt"
  "net/http"
  "time"
  
  "github.com/timcurless/lumberyard/Cassandra"
  "github.com/gocql/gocql"
)

func Post(w http.ResponseWriter, r *http.Request) {
  var errs []string
  var gocqlUUID gocql.UUID
  
  project, errs := FormToProject(r)
  
  var created = false
  
  // Add to Cassandra
  if len(errs) == 0 {
    fmt.Println("Creating a new Project")
    
    // generate a unique UUID for this project
		gocqlUUID = gocql.TimeUUID()
    nowTime := time.Now()
    
    if err := Cassandra.Session.Query(`
      INSERT INTO projects (id, name, email, updateTs, createdTs) VALUES (?, ?, ?, ?, ?)`,
			gocqlUUID, project.Name, project.Email, nowTime, nowTime).Exec(); err != nil {
			errs = append(errs, err.Error())
		} else {
			created = true
		}
    
  }
  
  // depending on whether we created the project, return the
	// resource ID in a JSON payload, or return our errors
	if created {
		fmt.Println("project_id", gocqlUUID)
		json.NewEncoder(w).Encode(NewProjectResponse{ID: gocqlUUID})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
