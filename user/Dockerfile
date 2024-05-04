# Step 1: Modules caching
FROM golang:1.22-alpine3.19 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22-alpine3.19 as builder
COPY --from=modules /go/pkg /go/pkg
COPY ./pkg /app/pkg
COPY ./user /app/user
COPY go.mod go.sum /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./user/cmd/server

# Step 3: Final
FROM scratch
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app", "--addr", "0.0.0.0:8081"]