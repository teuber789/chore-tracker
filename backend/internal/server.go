package internal

import (
	"context"
	"time"

	"github.com/teuber789/chore-tracker/internal/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TODO Use a database instead of in-memory maps
var children map[uint64]*gen.Child
var chores map[uint64]*gen.Chore
var completions map[uint64]*gen.ChoreCompletion

func init() {
	children = make(map[uint64]*gen.Child)
	chores = make(map[uint64]*gen.Chore)
	completions = make(map[uint64]*gen.ChoreCompletion)
}

func NewImplementedServer() *ImplementedChoreTrackerServer {
	return &ImplementedChoreTrackerServer{}
}

type ImplementedChoreTrackerServer struct {
	gen.UnimplementedChoreTrackerServer
}

func (ImplementedChoreTrackerServer) AddChild(ctx context.Context, req *gen.AddChildRequest) (*gen.Child, error) {
	nextId := uint64(0)
	for _, c := range children {
		if c.Id > nextId {
			nextId = c.Id
		}
	}
	nextId++

	child := &gen.Child{Id: nextId, Name: req.Name, Age: req.Age}
	children[nextId] = child
	return child, nil
}

func (*ImplementedChoreTrackerServer) CreateChore(ctx context.Context, req *gen.CreateChoreRequest) (*gen.Chore, error) {
	nextId := uint64(0)
	for _, c := range chores {
		if c.Id > nextId {
			nextId = c.Id
		}
	}
	nextId++

	chore := &gen.Chore{Id: nextId, Name: req.Name, Description: req.Description, Price: req.Price}
	chores[nextId] = chore
	return chore, nil
}

func (*ImplementedChoreTrackerServer) DeleteChore(ctx context.Context, req *gen.DeleteChoreRequest) (*emptypb.Empty, error) {
	delete(chores, req.ChoreId)
	return &emptypb.Empty{}, nil
}

func (*ImplementedChoreTrackerServer) GetChores(ctx context.Context, req *gen.GetChoresRequest) (*gen.GetChoresResponse, error) {
	payload := make([]*gen.Chore, len(chores))
	idx := 0
	for _, v := range chores {
		payload[idx] = v
		idx++
	}
	return &gen.GetChoresResponse{PageMetadata: &gen.PageMetadata{NextPageToken: ""}, Chores: payload}, nil
}

func (*ImplementedChoreTrackerServer) GetCompletedChores(context.Context, *gen.GetChoresRequest) (*gen.GetCompletedChoresResponse, error) {
	payload := make([]*gen.ChoreCompletion, len(completions))
	idx := 0
	for _, v := range completions {
		payload[idx] = v
		idx++
	}
	return &gen.GetCompletedChoresResponse{PageMetadata: &gen.PageMetadata{NextPageToken: ""}, ChoreCompletions: payload}, nil
}

func (*ImplementedChoreTrackerServer) MarkChoreCompleted(ctx context.Context, req *gen.MarkChoreCompletedRequest) (*emptypb.Empty, error) {
	nextId := uint64(0)
	for _, c := range completions {
		if c.Id > nextId {
			nextId = c.Id
		}
	}
	nextId++

	completed := &gen.ChoreCompletion{Id: nextId, ChoreId: req.ChoreId, ChildId: req.ChildId, CompletedTimestamp: uint32(time.Now().UnixMilli()), Paid: false}
	completions[nextId] = completed
	return &emptypb.Empty{}, nil
}
