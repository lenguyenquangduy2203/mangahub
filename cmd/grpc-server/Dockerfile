FROM golang:1.25-alpine AS builder
# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set up the working directory
WORKDIR /app/mangahub

# Copy root go.mod and go.sum (do not copy go.work into the build context)
# -- We don't want go build to use the workspace (go.work) during container builds
COPY go.mod go.sum* ./

# Copy all necessary files for the api-server while maintaining the directory structure
COPY cmd/grpc-server/main.go ./cmd/grpc-server/
COPY cmd/grpc-server/go.mod ./cmd/grpc-server/

# Copy proto files
COPY proto/manga/ ./proto/manga/

# Copy internal packages
COPY internal/grpc/ ./internal/grpc/
COPY internal/mangas/ ./internal/mangas/
COPY internal/users/ ./internal/users/

# Copy pkg packages
COPY pkg/errors/ ./pkg/errors/
COPY pkg/models/ ./pkg/models/
COPY pkg/pagination/ ./pkg/pagination/
COPY pkg/utils/ ./pkg/utils/

# Change to grpc-server directory for building
WORKDIR /app/mangahub/cmd/grpc-server
# Force module mode and ignore the workspace file at build-time
RUN GO111MODULE=on go build -mod=mod -o grpc-server .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/mangahub/cmd/grpc-server/grpc-server .
ENV PORT=50051
EXPOSE ${PORT}
ENTRYPOINT ["./grpc-server"]
