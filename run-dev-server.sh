#!/bin/bash

# Run Go server application.
go run settings.go handlers.go linkedin.go main.go &

# Run React application.
cd client
npm start &