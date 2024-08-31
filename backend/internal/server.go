package internal

import (
	"context"

	"github.com/teuber789/chore-tracker/internal/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewImplementedServer(store ChoreTrackerStore) *ImplementedChoreTrackerServer {
	return &ImplementedChoreTrackerServer{store: store}
}

type ImplementedChoreTrackerServer struct {
	gen.UnimplementedChoreTrackerServer
	store ChoreTrackerStore
}

func (i *ImplementedChoreTrackerServer) AddChild(ctx context.Context, req *gen.AddChildRequest) (*gen.Child, error) {
	return i.store.AddChild(ctx, req)
}

func (i *ImplementedChoreTrackerServer) CreateChore(ctx context.Context, req *gen.CreateChoreRequest) (*gen.Chore, error) {
	return i.store.CreateChore(ctx, req)
}

func (i *ImplementedChoreTrackerServer) DeleteChore(ctx context.Context, req *gen.DeleteChoreRequest) (*emptypb.Empty, error) {
	i.store.DeleteChore(ctx, req.ChoreId)
	return &emptypb.Empty{}, nil
}

func (i *ImplementedChoreTrackerServer) GetChores(ctx context.Context, req *gen.GetChoresRequest) (*gen.GetChoresResponse, error) {
	return i.store.GetChores(ctx, req)
}

func (i *ImplementedChoreTrackerServer) GetCompletedChores(ctx context.Context, req *gen.GetChoresRequest) (*gen.GetCompletedChoresResponse, error) {
	return i.store.GetCompletedChores(ctx, req)
}

func (i *ImplementedChoreTrackerServer) MarkChoreCompleted(ctx context.Context, req *gen.MarkChoreCompletedRequest) (*emptypb.Empty, error) {
	err := i.store.MarkChoreCompleted(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
