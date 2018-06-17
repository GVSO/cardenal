#! /bin/bash

# Run Go Server
glide install
go get github.com/pilu/fresh
rm -r $GOPATH/src
ln -s vendor $GOPATH/src
fresh &

# Run client
cd client
npm install
npm start