import os
import sys
import logging
import time
from heartbeat import HeartBeat
from execution_worker import ExecutionWorker
from receiver import Receiver

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
    receiver = Receiver(1)
    while True:
        ew = ExecutionWorker(receiver.model_path())
        ew.execute_packets(receiver.generate_data())

    heartbeat.join()
    logger.info('heartbeat terminated')
