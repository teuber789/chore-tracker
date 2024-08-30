# TODO IRL, these would be separate repos, so there wouldn't be hardcoded paths to subdirectories.
.PHONY: protos
protos:
	@protoc --go_out=backend/internal/gen --go_opt=paths=source_relative --go-grpc_out=backend/internal/gen --go-grpc_opt=paths=source_relative ./api/chore_tracker.proto --proto_path=./api	
