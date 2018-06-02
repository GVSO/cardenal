#!/bin/bash
trap 'killall' INT

killall() {
  trap '' INT TERM  # ignore INT and TERM while shutting down
  kill -TERM 0
  wait
}

# Run Go server application.
go run settings.go handlers.go linkedin.go main.go | sed "s/^/[Go Server] /" &

# Run React application.
cd client
npm start | sed "s/^/[React] /"  &

cat # wait forever