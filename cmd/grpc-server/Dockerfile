FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
# Only copy go.sum if it exists (ignore missing file)
RUN if [ -f go.sum ]; then cp go.sum .; fi
RUN go mod download
COPY . .
RUN go build -o grpc-server .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/grpc-server .
ENV PORT=50051
EXPOSE ${PORT}
ENTRYPOINT ["./grpc-server"]
