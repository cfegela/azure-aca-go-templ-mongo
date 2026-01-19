# Build stage
FROM golang:1-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy source code
COPY . .

# Generate templ templates
RUN /go/bin/templ generate

# Build both the server and seed binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./cmd/seed

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy both binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/seed .

# Copy static files
COPY --from=builder /app/web/static ./web/static

# Create startup script
RUN echo '#!/bin/sh' > /root/start.sh && \
    echo './seed || true' >> /root/start.sh && \
    echo './main' >> /root/start.sh && \
    chmod +x /root/start.sh

# Expose port
EXPOSE 8080

# Run the startup script
CMD ["/root/start.sh"]
