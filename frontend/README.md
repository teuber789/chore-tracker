# Chore Tracker Frontend

Frontend web app for the chore tracker. Still very rough. Based on GRPC Web's [Hello World example](https://github.com/grpc/grpc-web/blob/master/net/grpc/gateway/examples/helloworld/README.md?plain=1).

> ⚠️ The frontend only works with the GRPC service.

To run:

- Start the GRPC compose stack (see instructions in [the parent README](../README.md))
- In a new terminal, build the frontend for browser use and then serve locally:
  - `cd frontend`
  - `npm i`
  - `make dist`
  - `make serve`
  - Open Chrome to `http://localhost:3000`
  - Open your browser console and see the logs from the output.
