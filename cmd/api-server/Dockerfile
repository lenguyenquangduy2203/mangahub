FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set up the working directory
WORKDIR /app/mangahub

# Copy root go.mod and go.work files
COPY go.mod go.sum* ./
COPY go.work go.work.sum* ./

# Copy all necessary files for the api-server while maintaining the directory structure
COPY cmd/api-server/main.go ./cmd/api-server/
COPY cmd/api-server/go.mod ./cmd/api-server/

# Copy internal packages
COPY internal/api-server/ ./internal/api-server/
COPY internal/auth/ ./internal/auth/
COPY internal/socket/ ./internal/socket/

# Copy pkg packages
COPY pkg/errors/ ./pkg/errors/
COPY pkg/models/ ./pkg/models/
COPY pkg/pagination/ ./pkg/pagination/
COPY pkg/utils/ ./pkg/utils/

# Change to api-server directory for building
WORKDIR /app/mangahub/cmd/api-server
RUN go build -o api-server .

# Change to api-server directory for building
WORKDIR /app/mangahub/cmd/api-server

# Build the application
RUN go build -o api-server .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/mangahub/cmd/api-server/api-server .
ENV PORT=8080
EXPOSE ${PORT}
ENTRYPOINT ["./api-server"]
