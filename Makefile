.PHONY: grpc-down
grpc-down:
	docker compose -f infra/docker-compose-grpc.yaml down

.PHONY: grpc-up
grpc-up:
	docker compose -f infra/docker-compose-grpc.yaml up  --build --force-recreate

.PHONY: http-down
http-down:
	docker compose -f infra/docker-compose-http.yaml down

.PHONY: http-up
http-up:
	docker compose -f infra/docker-compose-http.yaml up --build --force-recreate

# TODO IRL, these would be separate repos, so there wouldn't be hardcoded paths to subdirectories.
.PHONY: protos
protos:
	@protoc --go_out=backend/internal/gen --go_opt=paths=source_relative --go-grpc_out=backend/internal/gen --go-grpc_opt=paths=source_relative ./api/chore_tracker.proto --proto_path=./api
	@protoc -I=./api ./api/chore_tracker.proto --js_out=import_style=commonjs:./frontend --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./frontend
	@cp ./frontend/*_pb.js ./load/runner
