# syntax=docker/dockerfile:1

##
## Build stage
##
FROM golang:1.18-buster as build

LABEL maintainer="Pavel Rezabek"

WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/
COPY api /app/api
COPY db /app/db
COPY cmd /app/cmd

RUN go mod download
RUN go build -o /go_challenge /app/cmd/go_challenge/

##
## Deploy stage
##
FROM gcr.io/distroless/base-debian10

WORKDIR /home/nonroot
COPY --from=build /go_challenge go_challenge

ARG PORT=8080
ARG HOST

ENV PORT=${PORT}
ENV HOST=${HOST}

EXPOSE ${PORT}
USER nonroot:nonroot

CMD [ "/home/nonroot/go_challenge" ]