import os
import sys
import logging
import atexit
from params_parser import parse_process_args
from heartbeat import HeartBeat


logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

stream = logging.StreamHandler(sys.stdout)
stream.setLevel(logging.INFO)
logger.addHandler(stream)


if __name__ == "__main__":

    process_args = parse_process_args()
    if process_args.get('error'):
        exit()

    def exit_handler():
        if heartbeat:
            heartbeat.gracefull_shutdown = True
        logger.info('[EXEC] Process shutdown successfully')
    atexit.register(exit_handler)

    logger.info(f'[EXEC]  Model executor process started with id-{os.getpid()}')
    heartbeat = HeartBeat(process_id=os.getpid())
    heartbeat.daemon = True
    heartbeat.start()

    logger.info('[EXEC] waiting for heartbeat to terimnate')
    heartbeat.join()
    logger.info('[EXEC] heartbeat terminated, main process going down')
