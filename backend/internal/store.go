package internal

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/teuber789/chore-tracker/internal/gen"
)

// IRL, creds would be injected as environment variables instead of hardcoded.
const (
	dbname   = "chore_tracker"
	host     = "127.0.0.1"
	port     = 5432
	user     = "chore-tracker-service"
	password = "choretrackerservicepassword"
)

func ConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
}

func NewChoreTrackerStore() (*store, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return &store{db: db}, nil
}

// IRL, transport layer structs wouldn't be part of the interface for the storage layer.
type ChoreTrackerStore interface {
	Close() error
	AddFamily(ctx context.Context, req *gen.AddFamilyRequest) (*gen.Family, error)
	AddChild(ctx context.Context, req *gen.AddChildRequest) (*gen.Child, error)
	CreateChore(ctx context.Context, req *gen.CreateChoreRequest) (*gen.Chore, error)
	DeleteChore(ctx context.Context, id uint64) error
	GetChores(ctx context.Context, req *gen.GetChoresRequest) (*gen.GetChoresResponse, error)
	GetCompletedChores(context.Context, *gen.GetChoresRequest) (*gen.GetCompletedChoresResponse, error)
	MarkChoreCompleted(ctx context.Context, req *gen.MarkChoreCompletedRequest) error
}

type store struct {
	db *sql.DB
}

func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) AddFamily(ctx context.Context, req *gen.AddFamilyRequest) (*gen.Family, error) {
	q := "INSERT INTO family (name) VALUES ($1) RETURNING id"
	id := uint64(0)
	err := s.db.QueryRow(q, req.Name).Scan(&id)
	if err != nil {
		return nil, err
	}

	f := &gen.Family{Id: id, Name: req.Name}
	return f, nil
}

func (s *store) AddChild(ctx context.Context, req *gen.AddChildRequest) (*gen.Child, error) {
	q := "INSERT INTO child (family_id, name, age) VALUES ($1, $2, $3) RETURNING id"
	id := uint64(0)
	err := s.db.QueryRow(q, req.FamilyId, req.Name, req.Age).Scan(&id)
	if err != nil {
		return nil, err
	}

	c := &gen.Child{Id: id, FamilyId: req.FamilyId, Name: req.Name, Age: req.Age}
	return c, nil
}

func (s *store) CreateChore(ctx context.Context, req *gen.CreateChoreRequest) (*gen.Chore, error) {
	q := "INSERT INTO chore (family_id, name, description, price) VALUES ($1, $2, $3, $4) RETURNING id"
	id := uint64(0)
	err := s.db.QueryRow(q, req.FamilyId, req.Name, req.Description, req.Price).Scan(&id)
	if err != nil {
		return nil, err
	}

	c := &gen.Chore{Id: id, FamilyId: req.FamilyId, Name: req.Name, Description: req.Description, Price: req.Price}
	return c, nil
}

func (s *store) DeleteChore(ctx context.Context, id uint64) error {
	q := "DELETE FROM chore WHERE id = $1"
	_, err := s.db.Exec(q, id)
	return err
}

func (s *store) GetChores(ctx context.Context, req *gen.GetChoresRequest) (*gen.GetChoresResponse, error) {
	// For this project, next page tokens are simply integers
	n := req.Pageable.PageSize
	page, err := strconv.Atoi(req.Pageable.PageToken)
	if err != nil {
		return nil, err
	}
	offset := uint32(page) * n

	// Execute query
	q := `SELECT id, name, description, price
		FROM chore
		WHERE family_id = $1
		ORDER BY created_at DESC
		OFFSET $2
		LIMIT $3`
	rows, err := s.db.Query(q, req.FamilyId, offset, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Read result set
	var arr []*gen.Chore = make([]*gen.Chore, 0, n)
	for rows.Next() {
		c := gen.Chore{}
		err = rows.Scan(&c.Id, &c.FamilyId, &c.Name, &c.Description, &c.Price)
		if err != nil {
			return nil, err
		}
		arr = append(arr, &c)
	}

	// Get the next token
	next := ""
	if uint32(len(arr)) == n {
		// We assume there are more records in the DB
		// This is good enough for this project
		next = strconv.Itoa(page + 1)
	}

	// Return
	res := &gen.GetChoresResponse{PageMetadata: &gen.PageMetadata{NextPageToken: next}, Chores: arr}
	return res, nil
}

func (s *store) GetCompletedChores(ctx context.Context, req *gen.GetChoresRequest) (*gen.GetCompletedChoresResponse, error) {
	// For this project, next page tokens are simply integers
	n := req.Pageable.PageSize
	page, err := strconv.Atoi(req.Pageable.PageToken)
	if err != nil {
		return nil, err
	}
	offset := uint32(page) * n

	// Execute query
	q := `SELECT id, family_id, child_id, chore_id, completed_timestamp, paid
		FROM chore_completion
		WHERE family_id = $1
		AND child_id = $2
		ORDER BY created_at DESC
		OFFSET $3
		LIMIT $4`
	rows, err := s.db.Query(q, req.FamilyId, req.ChildId, offset, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Read result set
	var arr []*gen.ChoreCompletion = make([]*gen.ChoreCompletion, 0, n)
	for rows.Next() {
		c := gen.ChoreCompletion{}
		err = rows.Scan(&c.Id, &c.FamilyId, &c.ChildId, &c.ChoreId, &c.CompletedTimestamp, &c.Paid)
		if err != nil {
			return nil, err
		}
		arr = append(arr, &c)
	}

	// Get the next token
	next := ""
	if uint32(len(arr)) == n {
		// We assume there are more records in the DB
		// This is good enough for this project
		next = strconv.Itoa(page + 1)
	}

	// Return
	res := &gen.GetCompletedChoresResponse{PageMetadata: &gen.PageMetadata{NextPageToken: next}, ChoreCompletions: arr}
	return res, nil
}

func (s *store) MarkChoreCompleted(ctx context.Context, req *gen.MarkChoreCompletedRequest) error {
	now := time.Now().UnixMilli()
	q := `UPDATE chore_completion
		SET completed_timestamp = $1
		WHERE family_id = $2
		AND child_id = $3
		AND chore_id = $4`
	_, err := s.db.Exec(q, now, req.FamilyId, req.ChildId, req.ChoreId)
	if err != nil {
		return err
	}

	return nil
}
