# Chore Tracker

A simple app that tracks the chores a child can do for their parents in order to get money.

# Project Structure

For convenience, this project is set up as a [monorepository](https://circleci.com/blog/monorepo-dev-practices/). In the real world, it would likely be split into several repos instead.

- [`backend`](./backend/) contains all of the code for the Golang-based GRPC web service.
- [`frontend`](./frontend/) contains all of the code for the React-based frontend.

# Prerequisites

- Make sure you have make installed
  - Run `brew install make`
- Install the latest Golang
- Protobufs
  - Run `brew install protobuf`
  - Run `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
  - Run `npm install -g protoc-gen-js`
- Install required gRPC compilers
  - Run `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
  - Run `brew install protoc-gen-grpc-web`
- Make sure you have docker installed
- Make sure you have npm installed
- Make sure you have python3 installed
  - Run `brew install python3`
- Make sure you have Chrome installed

The following are extremely useful for debugging, but not strictly necessary:

- grpcurl
  - Run `brew install grpcurl`
- curl
  - Run `brew install curl`

# Getting Started

> **Note:** This project exposes either a GRPC or an HTTP server for comparison purposes.

- Generate the grpc clients and their associated protobuf files:
  - Run `make protos`
- In a new terminal, run the docker compose stack:
  - `make compose-up`
- In a new terminal, run the backend you want:
  - `make [grpc|http]`. For example, to run the HTTP service, you would run `make http`.

> ⚠️ The frontend only works with the GRPC service. If you are using the HTTP service, there's no need to perform the next steps.

- In a new terminal, build the frontend for browser use and then serve locally:
  - `cd frontend`
  - `make dist`
  - `make serve`
  - Open Chrome to `http://localhost:3000`

# Future Development

These are the features I intentionally chose to ignore for the sake of this prototype. In the real world, these would obviously be addressed.

- Use a database instead of in-memory
- Real pagination (keyset instead of offset, with real next tokens)
- Authentication and authorization
- Permission-based access control (adults are allowed to create chores and children aren't, etc.)
- Logging
- Observability
- Metrics
- Testing
- Formatting / linting (and associated checks)
- Better dependency injection (wire?)
- Chores don't show up in the list of available chores after a child has completed them.
- Some chores can be performed more than once per day.
- Chore availability - certain chores might only be available on certain days, and therefore can only be completed on those days.
- Bug where all chores of a certain type are completed if multiple are performed
- Graceful HTTP server shutdown
- Intermediate representation for structs instead of reusing GRPC structs everywhere

# Work Log

- Evening of 29 Aug 2024: Added simple GRPC API for tracking chores. No DB, all storage is done in-memory.
- Evening of 30 Aug 2024:
  - Added frontend and envoy proxy; connected all to ensure they are working together.
  - Replaced in-memory storage with DB
- 31 Aug 2024:
  - Added HTTP server
  - Added multi-tenancy to make load testing easier (aka families)
