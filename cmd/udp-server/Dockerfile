FROM golang:1.25-alpine AS builder
WORKDIR /app/mangahub
COPY go.mod go.sum* ./
# Copy all necessary files for the udp-server while maintaining the directory structure
COPY cmd/udp-server/main.go ./cmd/udp-server/
COPY cmd/udp-server/go.mod ./cmd/udp-server/

# Copy internal packages
COPY internal/udp/ ./internal/udp/

# Change to udp-server directory for building
WORKDIR /app/mangahub/cmd/udp-server
RUN GO111MODULE=on go build -mod=mod -o udp-server .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/mangahub/cmd/udp-server/udp-server .
ENV PORT=9001
EXPOSE ${PORT}
ENTRYPOINT ["./udp-server"]
