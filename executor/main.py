import os
import sys
import logging
import atexit
from params_parser import parse_process_args
from heartbeat import HeartBeat
from receiver import Receiver
import time

BUSYFLAG = 0


def process():
    time.sleep(10)


logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

stream = logging.StreamHandler(sys.stdout)
stream.setLevel(logging.INFO)
logger.addHandler(stream)

fh = logging.FileHandler(f'logs/executor/logs-{os.getpid()}.log')
fh.setLevel(logging.INFO)
logger.addHandler(fh)

if __name__ == "__main__":
    args = parse_process_args()

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

    receiver = Receiver(videoPath=args.video_path,
                        modelPath=args.model_path,
                        modelConfigPath=args.model_config_path,
                        port=args.port)
    while True:
        receiver.get_batch_metadata()
        receiver.reply()
        heartbeat.curr_job_id = receiver.get_job_id()

        heartbeat.busy = True
        receiver.generate_data()
        process()
        heartbeat.busy = False

    heartbeat.join()
    logger.info('[EXEC] heartbeat terminated, main process going down')
