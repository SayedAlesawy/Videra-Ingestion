install:
	GO114MODULE=on go mod tidy

build: install
	GO114MODULE=on go build -o ingestion-engine.bin .

run: build
	./ingestion-engine.bin \
	-video-path=$(video-path) \
	-model-path=$(model-path) \
	-model-config-path=$(model-config-path) \
	-start-idx=$(start-idx) \
	-frame-count=$(frame-count)
