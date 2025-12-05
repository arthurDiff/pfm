# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/cli/main.go

# Run the application
run:
	@go run --race cmd/cli/main.go --port 8080 

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Remove unused dep
tidy:
	@go mod tidy