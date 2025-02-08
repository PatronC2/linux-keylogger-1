#!/bin/sh
set -e

# Start socat to create a pseudo-terminal and link it to /dev/tty0
SOCAT_OUTPUT=$(socat -d -d -ly PTY,raw,echo=0,link=/dev/tty0 2>&1) &
SOCAT_PID=$!

# Extract the actual PTY device name
TTY_DEVICE=$(echo "$SOCAT_OUTPUT" | grep -o '/dev/pts/[0-9]*' | head -n 1)

# Wait briefly to ensure /dev/tty0 is created
sleep 1

echo "Created pseudo-terminal: $TTY_DEVICE linked to /dev/tty0"

# Run tests
exec go test -v /app/keylogger
