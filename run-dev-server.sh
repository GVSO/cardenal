#!/bin/bash

# Run Go server application.
go run settings.go handlers.go linkedin.go main.go | sed "s/^/[Go Server] /" &

# Run React application.
cd client
npm start | sed "s/^/[React] /"  &

wait