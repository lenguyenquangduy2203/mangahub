FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
# Only copy go.sum if it exists (ignore missing file)
RUN if [ -f go.sum ]; then cp go.sum .; fi
RUN go mod download
COPY . .
RUN go build -o api-server .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/api-server .
ENV PORT=8080
EXPOSE ${PORT}
ENTRYPOINT ["./api-server"]
