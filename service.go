package lumberyard

import (
  "context"
  "errors"
  "fmt"
  "log"
  "strings"

  "github.com/gocql/gocql"
)

// Defines a deployment (i.e. a demo environment)
type Deployment struct {
  ID        gocql.UUID      `json:"id"`
  Name      string      `json:"name,omitempty"`
  Instances []string  `json:"instances,omitempty"`
}

// Main Service Interface
type Store interface {
  PostDeployment(ctx context.Context, d Deployment) error
  GetDeployments(ctx context.Context) ([]Deployment, error)
}

type CassandraStore struct {
  Session *gocql.Session
  Deployments []Deployment
}

type DataStoreFactory func() (Store)

func NewCassandraStore() (Store) {
  config := gocql.NewCluster("127.0.0.1")
  config.Keyspace = "lumberyard"
  session, err := config.CreateSession()
  if err != nil {
    panic(err)
  }
  return &CassandraStore {
    Session: session,
  }
}

var datastoreFactories = make(map[string]DataStoreFactory)

func Register(name string, factory DataStoreFactory) {
  if factory == nil {
    log.Panicf("Datastore factory %s does not exist.", name)
  }
  _, registered := datastoreFactories[name]
  if registered {
    log.Println("Datastore factory %s already registered. Ignoring.", name)
  }
  datastoreFactories[name] = factory
}

func init() {
  Register("cassandra", NewCassandraStore)
}

func CreateDatastore() (Store, error) {
  engineName := "cassandra"

  engineFactory, ok := datastoreFactories[engineName]
  if !ok {
    availableDatastores := make([]string, len(datastoreFactories))
    for k, _ := range datastoreFactories {
      availableDatastores = append(availableDatastores, k)
    }
    return nil, errors.New(fmt.Sprintf("Invalid datastore name. Must be one of: %s", strings.Join(availableDatastores, ", ")))
  }

  return engineFactory(), nil
}

// Main Services
func (s *CassandraStore) PostDeployment(ctx context.Context, d Deployment) error {

  if err := s.Session.Query(`
    INSERT INTO deployments (id, name, instances) VALUES(?, ?, ?)`,
    gocql.TimeUUID(), d.Name, d.Instances).Exec(); err != nil {
      return err
    } else {
      return nil
    }
}

func (s *CassandraStore) GetDeployments(ctx context.Context) ([]Deployment, error) {
  var deploymentList []Deployment
  m := map[string]interface{}{}

  iterable := s.Session.Query(`SELECT id,name,instances FROM deployments`).Iter()
  for iterable.MapScan(m) {
    deploymentList = append(deploymentList, Deployment{
      ID: m["id"].(gocql.UUID),
      Name: m["name"].(string),
      Instances: m["instances"].([]string),
    })
    m = map[string]interface{}{}
  }

  return deploymentList, nil
}
