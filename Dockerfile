FROM golang:1.22 as builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go mod tidy && \
    go install github.com/gordonklaus/ineffassign@latest && \
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o HalogenGhostCore .

FROM alpine
EXPOSE 1997
RUN apk add --no-cache tzdata
RUN mkdir /app /core
COPY --from=builder /app/HalogenGhostCore /app
CMD ["/app/HalogenGhostCore"]