FROM golang:1.23 AS builder
RUN mkdir /app
WORKDIR /app
COPY go.* .
RUN go mod download && \
    go install github.com/gordonklaus/ineffassign@latest && \
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
COPY . .
RUN ineffassign ./...
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o HalogenGhostCore .

FROM alpine
EXPOSE 1997
RUN apk add --no-cache tzdata
RUN mkdir /app /core
COPY --from=builder /app/HalogenGhostCore /app
CMD ["/app/HalogenGhostCore"]