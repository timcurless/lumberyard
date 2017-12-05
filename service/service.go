package service

import (
	"context"
	"errors"
	"sync"

	"github.com/gocql/gocql"
)

// Service representing Lumberyard state
type Service interface {
	PostProject(ctx context.Context, p Project) (string, error)
	GetProject(ctx context.Context, id string) (Project, error)
}

// Project is a top level Project resource
type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	UpdateTs  string `json:"update-ts"`
	CreatedTs string `json:"created-ts"`
}

var (
	// ErrAlreadyExists is an error used when something already exists (create, not overwrite)
	ErrAlreadyExists = errors.New("already exists")
	// ErrNotFound is an error used when something is not found
	ErrNotFound = errors.New("not found")
)

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]Project
}

type cassandraService struct {
	db *gocql.Session
}

// NewInmemService creates an in-memory database
func NewInmemService() Service {
	return &inmemService{
		m: map[string]Project{},
	}
}

// NewCassandraService creates a Service persisting in Cassadndra
func NewCassandraService(uri, username, password, keyspace string) Service {
	cluster := gocql.NewCluster(uri)
	cluster.Keyspace = keyspace
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		panic("Failed initializing Cassandra Store")
	}
	return &cassandraService{
		db: session,
	}
}

func (s *inmemService) PostProject(ctx context.Context, p Project) (string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[p.ID]; ok {
		return "", ErrAlreadyExists
	}
	s.m[p.ID] = p
	return p.ID, nil
}

func (s *inmemService) GetProject(ctx context.Context, id string) (Project, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[id]
	if !ok {
		return Project{}, ErrNotFound
	}
	return p, nil
}

func (s *cassandraService) PostProject(ctx context.Context, p Project) (string, error) {

	err := s.db.Query(`INSERT INTO projects (id, name, email, update_ts, created_ts) VALUES (?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Email, p.UpdateTs, p.CreatedTs).Exec()
	if err != nil {
		return "", err
	}
	return p.ID, nil
}

func (s *cassandraService) GetProject(ctx context.Context, id string) (Project, error) {
	var p Project
	var found = false
	m := map[string]interface{}{}

	query := "SELECT id,name,email,update_ts,created_ts FROM projects WHERE id=? LIMIT 1"
	iterable := s.db.Query(query, id).Consistency(gocql.One).Iter()

	for iterable.MapScan(m) {
		found = true
		p = Project{
			ID:        m["id"].(string),
			Name:      m["name"].(string),
			Email:     m["email"].(string),
			UpdateTs:  m["update_ts"].(string),
			CreatedTs: m["created_ts"].(string),
		}
	}

	if !found {
		return Project{}, ErrNotFound
	}
	return p, nil

}
