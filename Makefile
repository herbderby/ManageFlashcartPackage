BINARY = ext/server/flashcart-tools
MCPB = flashcart-tools.mcpb
LDFLAGS = -s -w
SOURCES = main.go volumes.go filesystem.go bytes.go network.go \
          archive.go image.go json_tools.go skill.go hash.go \
          nointro.go nointro.json.gz go.mod go.sum

.PHONY: build pack clean vet test gen-nointro

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $@

pack: $(BINARY)
	cd ext && zip -r ../$(MCPB) manifest.json server/flashcart-tools

vet:
	go vet ./...

test:
	go test ./...

gen-nointro:
	go run ./tools

clean:
	rm -f $(BINARY) $(MCPB)
