BINARY_DIR = ext/server
LAUNCHER = $(BINARY_DIR)/flashcart-tools
DARWIN_ARM64 = $(BINARY_DIR)/flashcart-tools-darwin-arm64
LINUX_ARM64 = $(BINARY_DIR)/flashcart-tools-linux-arm64
MCPB = flashcart-tools.mcpb
LDFLAGS = -s -w
SOURCES = main.go volumes.go volumes_darwin.go volumes_linux.go \
          filesystem.go bytes.go network.go \
          archive.go image.go json_tools.go skill.go hash.go \
          nointro.go nointro.json.gz go.mod go.sum

.PHONY: build pack clean vet test gen-nointro

build: $(DARWIN_ARM64) $(LINUX_ARM64)

$(DARWIN_ARM64): $(SOURCES)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $@

$(LINUX_ARM64): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $@

pack: $(DARWIN_ARM64) $(LINUX_ARM64)
	cd ext && zip -r ../$(MCPB) manifest.json \
		server/flashcart-tools \
		server/flashcart-tools-darwin-arm64 \
		server/flashcart-tools-linux-arm64

vet:
	go vet ./...

test:
	go test ./...

gen-nointro:
	go run ./tools

clean:
	rm -f $(DARWIN_ARM64) $(LINUX_ARM64) $(MCPB)
