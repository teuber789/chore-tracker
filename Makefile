.PHONY: envoy
envoy:
	@$ docker run --rm -v $(shell pwd)/infra/envoy.yaml:/etc/envoy/envoy.yaml:ro -p 8080:8080 envoyproxy/envoy:v1.22.0

# TODO IRL, these would be separate repos, so there wouldn't be hardcoded paths to subdirectories.
.PHONY: protos
protos:
	@protoc --go_out=backend/internal/gen --go_opt=paths=source_relative --go-grpc_out=backend/internal/gen --go-grpc_opt=paths=source_relative ./api/chore_tracker.proto --proto_path=./api
	@protoc -I=./api ./api/chore_tracker.proto --js_out=import_style=commonjs:./frontend --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./frontend
