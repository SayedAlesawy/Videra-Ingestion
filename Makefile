install:
	GO114MODULE=on go mod tidy

build: install
	GO114MODULE=on go build -o ingestion-engine.bin .

run: build
	./ingestion-engine.bin \
	-execution-group-id=$(execution-group-id) \
	-video-path=$(video-path) \
	-model-path=$(model-path) \
	-model-config-path=$(model-config-path) \
	-start-idx=$(start-idx) \
	-frame-count=$(frame-count)

# Example command
# make execution-group-id=exec_group_1 video-path=/home/sayed/series/test/file.txt model-path=/home/sayed/series/test/file.txt model-config-path=/home/sayed/series/test/file.txt start-idx=5 frame-count=50 run
