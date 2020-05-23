install:
	GO114MODULE=on go mod tidy

build: install
	GO114MODULE=on go build -o ingestion-engine.bin .

run: build
	./ingestion-engine.bin
