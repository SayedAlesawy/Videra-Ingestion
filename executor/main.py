import os
import sys
import logging
import atexit
from params_parser import parse_process_args
from heartbeat import HeartBeat
from receiver import Receiver
from execution_worker import ExecutionWorker


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
    pid = os.getpid()

    logger.info(f'[EXEC]  Model executor process started with id-{pid}')
    heartbeat = HeartBeat(process_id=pid)
    heartbeat.send_heartbeat()
    heartbeat.daemon = True
    heartbeat.start()

    logger.info('[EXEC] waiting for heartbeat to terimnate')

    receiver = Receiver(cache_prefix=args.execution_group_id, pid=pid)
    executor = ExecutionWorker(args.model_path, args.video_path)
    while True:
        job_meta, job_key = receiver.get_new_job()
        try:
            heartbeat.curr_job_id = job_key.decode('utf-8')
        except Exception as e:
            logger.exception(f'[EXEC] job key is malformed : {e}')
            continue

        heartbeat.busy = True
        heartbeat.send_heartbeat()

        executor.execute(job_meta)

        receiver.mark_job_as_done(job_key)
        heartbeat.busy = False
        heartbeat.send_heartbeat()

    heartbeat.join()
    logger.info('[EXEC] heartbeat terminated, main process going down')
