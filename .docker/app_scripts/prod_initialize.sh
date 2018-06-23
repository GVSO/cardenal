#!/bin/bash

# Build client
cd client
npm install --production

# Build Go Server
cd ..
dep ensure
go build -o cardenal
./cardenal
