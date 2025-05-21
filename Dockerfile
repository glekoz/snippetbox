FROM golang:1.24.2-bookworm

WORKDIR /snippetbox

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .