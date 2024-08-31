.PHONY: compose-down
compose-down:
	docker compose -f infra/docker-compose.yaml down

.PHONY: compose-up
compose-up:
	docker compose -f infra/docker-compose.yaml up

.PHONY: grpc
grpc:
	cd backend && go run main.go -server=grpc

.PHONY: http
http:
	cd backend && go run main.go -server=http

# TODO IRL, these would be separate repos, so there wouldn't be hardcoded paths to subdirectories.
.PHONY: protos
protos:
	@protoc --go_out=backend/internal/gen --go_opt=paths=source_relative --go-grpc_out=backend/internal/gen --go-grpc_opt=paths=source_relative ./api/chore_tracker.proto --proto_path=./api
	@protoc -I=./api ./api/chore_tracker.proto --js_out=import_style=commonjs:./frontend --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./frontend
