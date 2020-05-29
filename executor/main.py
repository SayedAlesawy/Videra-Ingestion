import os
import sys
import logging
import time
from heartbeat import HeartBeat
from params_parser import parse_process_args
# from execution_worker import ExecutionWorker
from receiver import Receiver

BUSYFLAG = 0


def process():
    pass


logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

stream = logging.StreamHandler(sys.stdout)
stream.setLevel(logging.INFO)
logger.addHandler(stream)

if __name__ == "__main__":
    logger.info(f'Model executor process started with id-{os.getpid()}')
    heartbeat = HeartBeat()
    heartbeat.daemon = False
    heartbeat.start()
    logger.info('waiting for heartbeat to terimnate')
    time.sleep(6)
    # heartbeat.gracefull_shutdown = True

    args = parse_process_args()
    receiver = Receiver(videoPath=args.video_path, modelPath=args.model_path, modelConfigPath=args.model_config_path)

    while True:
        receiver.get_batch_metadata()
        BUSYFLAG = 1
        process()
        BUSYFLAG = 0

    heartbeat.join()
    logger.info('heartbeat terminated')
