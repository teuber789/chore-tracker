# Load Testing

This directory contains a Golang program that can be used to load testing the Chore Tracker application.

# Running Load Tests

To run a load test:

- Using the terminal, stand up the stack for the server you want to load test (either GRPC or HTTP)
- In a new terminal, install all prerequisites:
    - `cd load`
    - `make prereqs`
- Run the load test with the parameters you choose:

```
go run main.go -server <grpc|http> -host 192.168.1.229 -seconds 3 -users 3
```

- Arguments are:
  - `server`: Whether the GRPC or the HTTP load test should be run. Required.
  - `host`: The host that the server is running on. Defaults to localhost if not specified.
  - `seconds`: The number of seconds the test should run before terminating. Defaults to 300 (5 minutes).
  - `users` The number of concurrent users the test should simulate. Defaults to 1.

# Notes:

- No load testing client I am aware of supports GRPC Web's protocol. As such, I could not use any of them to compare GRPC Web performance to that of a RESTful API. I wrote this load test runner to fill this gap.
