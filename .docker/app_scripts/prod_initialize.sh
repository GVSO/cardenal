#! /bin/bash

# Build client
cd client
npm install --production

# Build Go Server
cd ..
glide install
rm -r $GOPATH/src
ln -s vendor $GOPATH/src
go build -o cardenal
./cardenal
