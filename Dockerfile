FROM golang:1.24.2-alpine AS base

RUN apk add --no-cache \
  curl \
  git \
  tzdata \
  bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM base AS dev

RUN apk add --no-cache make
