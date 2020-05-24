import os
import sys
import logging
import atexit
from heartbeat import HeartBeat


logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

stream = logging.StreamHandler(sys.stdout)
stream.setLevel(logging.INFO)
logger.addHandler(stream)

if __name__ == "__main__":

    def exit_handler():
        if heartbeat:
            heartbeat.gracefull_shutdown = True
        print('[EXEC] Process Tear Down')
    atexit.register(exit_handler)

    logger.info(f'[EXEC]  Model executor process started with id-{os.getpid()}')
    heartbeat = HeartBeat()
    heartbeat.daemon = True
    heartbeat.start()
    logger.info('[EXEC] waiting for heartbeat to terimnate')
    # heartbeat.gracefull_shutdown = True
    heartbeat.join()
    logger.info('[EXEC] heartbeat terminated, main process going down')
