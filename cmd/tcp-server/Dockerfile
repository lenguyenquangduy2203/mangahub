FROM golang:1.25-alpine AS builder
# Set up the working directory
WORKDIR /app/mangahub

# Copy root go.mod for the root module
COPY go.mod go.sum* ./

# Copy all necessary files for the tcp-server while maintaining the directory structure
COPY cmd/tcp-server/main.go ./cmd/tcp-server/
COPY cmd/tcp-server/go.mod ./cmd/tcp-server/

# Copy internal packages
COPY internal/tcp/ ./internal/tcp/

# Change to tcp-server directory for building
WORKDIR /app/mangahub/cmd/tcp-server
RUN GO111MODULE=on go build -mod=mod -o tcp-server .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/mangahub/cmd/tcp-server/tcp-server .
ENV PORT=9000
EXPOSE ${PORT}
ENTRYPOINT ["./tcp-server"]
