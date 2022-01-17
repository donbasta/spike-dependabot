FROM golang:alpine AS builder

LABEL maintainer="Cloud Automation <cloud-automation@gopay.co.id>"

WORKDIR /app
ARG GITCONFIG
ADD . /app
RUN apk add build-base gcc git
ENV GOPRIVATE=source.golabs.io
RUN echo $GITCONFIG | base64 -d > ~/.gitconfig
RUN cd /app
RUN go mod vendor && CGO_ENABLED=1 go build -mod=vendor -race -a  -o bin/scp-dependency-manager main.go

FROM alpine

ARG GITCONFIG
ENV GOPRIVATE=source.golabs.io
RUN echo $GITCONFIG | base64 -d > ~/.gitconfig

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /app/bin/scp-dependency-manager /app
COPY --from=builder /app/migrations /app/migrations
