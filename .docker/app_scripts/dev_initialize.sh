#!/bin/bash

# Run Go Server
dep ensure
fresh &

# Run client
cd client
npm install
npm start