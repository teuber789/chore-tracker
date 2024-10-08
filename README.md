# Chore Tracker

A simple app that tracks the chores a child can do for their parents in order to get money.

# About

This project has two purposes:

- Showcase my skills with various technologies, including Golang backend development, GRPC, understanding of DDD, understanding of API best practices, etc.
- Provide a playground to experiment with [GRPC Web](https://github.com/grpc/grpc-web). Prior to this project, I had never utilized it before and wanted to play around with it.

That being said, this project is still a work in progress and has some rough edges (see the [future development section](#future-development)). It is not meant to be perfect, so much as to give a glimpse of what I am capable of.

> Throughout the repository, I frequently make comments starting with `IRL`. These comments explain how something would be different if I were deploying the code to a production environment in real life.

# Project Structure

For convenience, this project is set up as a [monorepository](https://circleci.com/blog/monorepo-dev-practices/). In the real world, it would likely be split into several repos instead.

- [`api`](./api/) contains all of the `.proto` files necessary for generating the GPRC API.
- [`backend`](./backend/) contains all of the code for the Golang-based GRPC web service.
- [`frontend`](./frontend/) contains all of the code for the React-based frontend.
- [`infra`](./infra/) contains all infrastructure-related code and configuration.
- [`load`](./load/) contains a small application to perform load testing.

# Prerequisites

- Make sure you have make installed
  - Run `brew install make`
- Make sure you have node installed
- Make sure you have docker installed
- Make sure you have python3 installed
  - Run `brew install python3`
- Make sure you have the latest Golang installed
- Protobufs
  - Run `brew install protobuf`
  - Run `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
  - Run `npm install -g protoc-gen-js`
- Install required gRPC compilers
  - Run `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
  - Run `brew install protoc-gen-grpc-web`
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
- In a new terminal, run the docker compose for the backend service you want:
  - `make [grpc|http]-up`. For example, to run the HTTP service, you would run `make http-up`.
- When you're done and want to tear down the stack, simply run `make [grpc|http]-down`. For example, if you were running the HTTP server and are now done, run `make http-down`.

To run the [frontend](./frontend/README.md) or the [load test application](./load/README.md), please see the READMEs in their respective directories.

# Future Development

These are areas I know this project can be improved upon.

- The frontend webapp needs to be built out. Currently, it only runs an example app.
- The data model isn't very flexible and doesn't allow for a variety of important use cases. It also creates a number of bugs. I recognize these but haven't fixed them yet.
  - Chores shouldn't show up in the list of available chores after a child has completed them.
  - Some chores can be performed more than once per day.
  - Chore availability - certain chores might only be available on certain days, and therefore can only be completed on those days.
  - Bug where all chores of a certain type are completed if multiple are performed
- Real pagination (keyset instead of offset, with real next tokens)
- Authentication and authorization
- Permission-based access control (adults are allowed to create chores and children aren't, etc.)
- Logging
- Observability
- Metrics
- Testing
- Formatting / linting (and associated checks)
- Better dependency injection (wire?)
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
  - Added load test module. WIP; executes a long-running process and terminates it when the context times out.
- 2 Sep 2024:
  - Added load test script for GRPC
  - Added metrics to load test script
  - Added load tests cript for HTTP (with metrics)
  - Changed the server to run in a Docker container. This will make it easier to load test.
  - Added a host parameter to the load test binary. This makes possible to test against a service running on another machine.
- 3 Sep 2024: Clarified purpose of repository and instructions.
