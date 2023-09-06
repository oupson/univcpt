FROM golang:1.21

# Set destination for COPY
WORKDIR /app

ARG DEFAULT_ICAL_URL
ENV ICAL_URL=$DEFAULT_ICAL_URL

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY cmd ./cmd
COPY internal ./internal

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /univcpt ./cmd/univcpt/main.go

# Expose port
EXPOSE 8000

# Run
CMD ["/univcpt"]