default: build

run:
	go run ./cmd/bolt

build:
	go build ./cmd/bolt

install:
	go install ./cmd/bolt
