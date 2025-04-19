.PHONY: all build clean

APP := snipcode
GOOS_ARCH := linux/amd64 linux/arm64 linux/386 linux/arm darwin/amd64 darwin/arm64 windows/amd64 windows/arm64 windows/386
DIST_DIR := dist

all: build

build:
	@echo "Building binaries..."
	@mkdir -p $(DIST_DIR)
	@for t in $(GOOS_ARCH); do \
		os=$${t%/*}; arch=$${t#*/}; \
		bin=$(DIST_DIR)/$${APP}-$${os}-$${arch}; \
		if [ "$$os" = "windows" ]; then bin:=$${bin}.exe; fi; \
		echo "  Building for $$os/$$arch..."; \
		GOOS=$$os GOARCH=$$arch go build -ldflags="-s -w" -o $$bin .; \
	done
	@echo "Build complete. Binaries in $(DIST_DIR)/"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(DIST_DIR)
	@echo "Clean complete."
