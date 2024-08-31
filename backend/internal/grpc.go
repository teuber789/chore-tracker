package internal

import (
	"context"

	"github.com/teuber789/chore-tracker/internal/gen"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewGrpcServer(store ChoreTrackerStore) *grpcSrv {
	return &grpcSrv{store: store}
}

type grpcSrv struct {
	gen.UnimplementedChoreTrackerServer
	store ChoreTrackerStore
}

func (s *grpcSrv) AddChild(ctx context.Context, req *gen.AddChildRequest) (*gen.Child, error) {
	return s.store.AddChild(ctx, req)
}

func (s *grpcSrv) CreateChore(ctx context.Context, req *gen.CreateChoreRequest) (*gen.Chore, error) {
	return s.store.CreateChore(ctx, req)
}

func (s *grpcSrv) DeleteChore(ctx context.Context, req *gen.DeleteChoreRequest) (*emptypb.Empty, error) {
	s.store.DeleteChore(ctx, req.ChoreId)
	return &emptypb.Empty{}, nil
}

func (s *grpcSrv) GetChores(ctx context.Context, req *gen.GetChoresRequest) (*gen.GetChoresResponse, error) {
	return s.store.GetChores(ctx, req)
}

func (s *grpcSrv) GetCompletedChores(ctx context.Context, req *gen.GetChoresRequest) (*gen.GetCompletedChoresResponse, error) {
	return s.store.GetCompletedChores(ctx, req)
}

func (s *grpcSrv) MarkChoreCompleted(ctx context.Context, req *gen.MarkChoreCompletedRequest) (*emptypb.Empty, error) {
	err := s.store.MarkChoreCompleted(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
