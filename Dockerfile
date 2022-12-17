FROM golang:1.19-alpine AS builder
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o api-logger
RUN mv /build/api-logger /

FROM alpine
WORKDIR /
COPY --from=builder /api-logger /api-logger
CMD ["/api-logger"]