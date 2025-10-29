FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
# Only copy go.sum if it exists (ignore missing file)
RUN if [ -f go.sum ]; then cp go.sum .; fi
RUN go mod download
COPY . .
RUN go build -o udp-server .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/udp-server .
ENV PORT=9001
EXPOSE ${PORT}/udp
ENTRYPOINT ["./udp-server"]
