FROM golang:stretch

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

# Build directories
RUN mkdir -p /go/src/github.com/gvso/cardenal
COPY . /go/src/github.com/gvso/cardenal

# Go dep!
RUN go get -u github.com/golang/dep/cmd/dep

# Installs Node.js
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash - && \
  apt-get install -y nodejs

CMD ["./.docker/app_scripts/prod_initialize.sh"]