# First stage: Run tests with TTY simulation
FROM golang:1.23.6 AS test

WORKDIR /app
COPY . ./
RUN go mod download

# Install socat to simulate /dev/tty0
RUN apt-get update && apt-get install -y socat

# Copy entrypoint script to start socat before running tests
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

# Second stage: Build the binary
FROM golang:1.23.6 AS build

WORKDIR /app
COPY --from=test /app ./
RUN go build -o keylogger ./main.go

# Third stage: Copy the binary for host output
FROM alpine:latest AS output
WORKDIR /output
COPY --from=build /app/keylogger .
