package service

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/gocql/gocql"
)

// Service representing Lumberyard state
type Service interface {
	PostProject(ctx context.Context, p Project) (string, error)
	GetProject(ctx context.Context, id string) (Project, error)
	PostStack(ctx context.Context, projectID string, s Stack) (string, error)
}

// Project is a top level Project resource
type Project struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	UpdateTs  string  `json:"update-ts"`
	CreatedTs string  `json:"created-ts"`
	Stacks    []Stack `json:"stacks,omitempty"`
}

// Stack is a struct representing a collection of assets for a project
type Stack struct {
	ID        string  `json:"id"`
	Assets    []Asset `json:"assets,omitempty"`
	UpdateTs  string  `json:"update-ts"`
	CreatedTs string  `json:"created-ts"`
}

// Asset is a struct representing an asset (i.e EC2 Instance or Load Balancer)
type Asset struct {
	ID   string `json:"id"`
	Type string `json:"type"`
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
func NewCassandraService(uri string) Service {
	cluster := gocql.NewCluster(uri)
	cluster.Keyspace = "system"
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		panic("Failed initializing Cassandra master Session")
	}

	if err := initSchema(session); err != nil {
		panic("Failed initializing Cassandra Schema: " + err.Error())
	}

	cluster.Keyspace = "lumberyard"
	mainSession, err := cluster.CreateSession()
	// Find a way to defer/close session

	if err != nil {
		panic("Failed initializing Cassandra main Session")
	}

	return &cassandraService{
		db: mainSession,
	}
}

func initSchema(s *gocql.Session) error {
	if err := s.Query(`CREATE KEYSPACE IF NOT EXISTS lumberyard
										 WITH replication = {
											 'class' : 'SimpleStrategy',
											 'replication_factor' : 1
										 }`).Exec(); err != nil {
		return err
	}

	/*
		if err := s.Query(`CREATE TYPE IF NOT EXISTS lumberyard.stack (
			 									 id         text,
												 assets     map<text, text>,
												 update_ts  text,
												 created_ts text
											 )`).Exec(); err != nil {
			return err
		}*/

	if err := s.Query(`CREATE TABLE IF NOT EXISTS lumberyard.projects (
											 id         text,
											 name       text,
											 email      text,
											 update_ts  text,
											 created_ts text,
											 stacks     text,
											 PRIMARY KEY (id)
										)`).Exec(); err != nil {
		return err
	}

	return nil
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

func (s *inmemService) PostStack(ctx context.Context, projectID string, st Stack) (string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	p, ok := s.m[projectID]
	if !ok {
		return "", ErrNotFound
	}
	for _, stack := range p.Stacks {
		if stack.ID == st.ID {
			return "", ErrAlreadyExists
		}
	}
	p.Stacks = append(p.Stacks, st)
	s.m[projectID] = p
	return p.ID, nil
}

func (s *cassandraService) PostProject(ctx context.Context, p Project) (string, error) {
	json, _ := json.Marshal(p.Stacks)
	err := s.db.Query(`INSERT INTO projects (id, name, email, update_ts, created_ts, stacks) VALUES (?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Email, p.UpdateTs, p.CreatedTs, json).Exec()
	if err != nil {
		return "", err
	}
	return p.ID, nil
}

func (s *cassandraService) GetProject(ctx context.Context, id string) (Project, error) {
	var p Project
	var found = false
	m := map[string]interface{}{}

	query := "SELECT id, name, email, update_ts, created_ts, stacks FROM projects WHERE id=? LIMIT 1"
	iterable := s.db.Query(query, id).Consistency(gocql.One).Iter()

	for iterable.MapScan(m) {
		found = true
		var stacku []Stack
		if err := fromJson(m["stacks"].(string), &stacku); err != nil {
			return Project{}, err
		}
		p = Project{
			ID:        m["id"].(string),
			Name:      m["name"].(string),
			Email:     m["email"].(string),
			UpdateTs:  m["update_ts"].(string),
			CreatedTs: m["created_ts"].(string),
			Stacks:    stacku,
		}

	}

	if !found {
		return Project{}, ErrNotFound
	}
	return p, nil

}

func (s *cassandraService) PostStack(ctx context.Context, projectID string, st Stack) (string, error) {

	query := "UPDATE projects SET stacks = ? + stacks WHERE id = ?"

	err := s.db.Query(query, st, projectID).Exec()
	if err != nil {
		return "", err
	}
	return st.ID, nil
}

func fromJson(jsonSrc string, s *[]Stack) error {
	return json.Unmarshal([]byte(jsonSrc), s)
}
