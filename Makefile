# Variables
BINARY_NAME=keylogger
OUTPUT_DIR=./bin
DOCKER_IMAGE=patron-linux-keylogger-1
TEST_CONTAINER=test-container

.PHONY: test-build test build extract clean

# Build the test image with socat for TTY simulation
test-build:
	docker build --target=test -t $(DOCKER_IMAGE) .

# Run the test container, ensuring socat creates /dev/tty0 at runtime
test:
	docker run --rm --cap-add=SYS_ADMIN --device /dev/ptmx --name $(TEST_CONTAINER) $(DOCKER_IMAGE)

# Build the Go binary inside the container
build:
	docker build --target=build -t $(DOCKER_IMAGE) .

# Extract the built binary from the container to the host system
extract:
	mkdir -p $(OUTPUT_DIR)
	docker create --name $(TEST_CONTAINER) $(DOCKER_IMAGE)
	docker cp $(TEST_CONTAINER):/app/$(BINARY_NAME) $(OUTPUT_DIR)/
	docker rm $(TEST_CONTAINER)
	chmod +x $(OUTPUT_DIR)/$(BINARY_NAME)

# Run all steps: test-build, test, build, extract
all: test-build test build extract

# Clean up Docker images and built binary
clean:
	docker rmi -f $(DOCKER_IMAGE) || true
	rm -rf $(OUTPUT_DIR)
