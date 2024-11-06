FROM golang:1.23-alpine AS based_builder

WORKDIR /src/go
COPY go.mod go.sum ./
RUN go mod download

#==============================================================================
FROM based_builder AS builder

ARG SERVICE

WORKDIR /src/go
COPY . ./
RUN go build -o /app "./${SERVICE}"

#==============================================================================
FROM alpine:latest

COPY --from=builder /app ./
ENTRYPOINT [ "./app" ]
