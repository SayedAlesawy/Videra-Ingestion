import os
import sys
import logging
import time
from heartbeat import HeartBeat
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

    modelPath = ""
    videoPath = ""
    modelConfigPath = ""
    receiver = Receiver(videoPath, modelPath, modelConfigPath)

    while True:
        receiver.get_batch_metadata()
        BUSYFLAG = 1
        process()
        BUSYFLAG = 0

    heartbeat.join()
    logger.info('heartbeat terminated')
