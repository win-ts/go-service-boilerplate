FROM golang:1.23-alpine AS based_builder

WORKDIR /src/go
COPY go.mod go.sum ./
RUN go mod download

#==============================================================================
FROM based_builder AS builder

WORKDIR /src/go
COPY . ./
RUN go build -o /app .

#==============================================================================
FROM alpine:latest

ENV TZ=Asia/Bangkok
COPY --from=builder /app ./
ENTRYPOINT [ "./app" ]
