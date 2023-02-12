.PHONY = test clean build cross

test:
	go test ./pkg/psu

clean:
	rm -rf ./build 2>/dev/null || true
	rm -rf ./fyne-cross 2>/dev/null || true
	rm Icon.png 2>/dev/null || true

build:
	go build -ldflags="-s -w" -o build/gui ./cmd/gui
	cp cmd/gui/config.json build/config.json

os ?= windows
cross:
	fyne-cross ${os} ./cmd/gui -release
