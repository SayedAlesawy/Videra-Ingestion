install:
	GO111MODULE=on go mod tidy

build: install
	GO111MODULE=on go build -o ingestion-engine.bin .

run: build
	./ingestion-engine.bin
