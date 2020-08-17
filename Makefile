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
	-code-path=$(code-path) \
	-video-token=$(video-token) \
	-start-idx=$(start-idx) \
	-frame-count=$(frame-count)

# Example command
# orch: make execution-group-id=exec_group_1 video-path=/home/sayed/series/test/file.txt model-path=/home/sayed/series/test/file.txt model-config-path=/home/sayed/series/test/file.txt start-idx=5 frame-count=50 -code-path ./path/to/code_driectory -video-token=mm.mp4 run
# single process: python ./actions/main.py  -execution-group-id=exec_group_1 -video-path=./mm.mp4 -model-path=./model_config/t.pickle -model-config-path=./model_config/model_config.json -code-path=./actions/executor/ -video-token=mm.mp4