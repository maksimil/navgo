run:
	go run ./cmd/navgo

build:
	go build -ldflags "-s -w" ./cmd/navgo
