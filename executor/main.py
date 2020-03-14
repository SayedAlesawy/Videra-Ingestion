import os
import sys
import logging
import time
from heartbeat import HeartBeat

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
    heartbeat.join()
    logger.info('heartbeat terminated')
