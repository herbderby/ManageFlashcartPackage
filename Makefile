BINARY = HelloSDCard/bin/hello-sdcard-darwin-arm64
LDFLAGS = -s -w

.PHONY: build clean vet test

build: $(BINARY)

$(BINARY): main.go volumes.go go.mod go.sum
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $@

vet:
	go vet ./...

test:
	go test ./...

clean:
	rm -f $(BINARY)
