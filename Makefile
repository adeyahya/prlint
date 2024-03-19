# Build the Docker image
.PHONY: build
build:
	docker build -t prlint-image --target executable .
	docker build -t prlint-image-test --target test .

# Run the Docker container
.PHONY: test
test:
	docker run --rm -e TITLE="fix: glob pattern" prlint-image-test -s "be92a1dee5e3b7b318fc021dcfd98d33fbc7c8e3"

# Clean up Docker artifacts
.PHONY: clean
clean:
	docker image rm prlint-image prlint-image-test

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build             - Build Docker image"
	@echo "  make test              - Run Test Docker container"
	@echo "  make clean             - Clean Docker artifacts"
	@echo "  make help              - Show this help message"
