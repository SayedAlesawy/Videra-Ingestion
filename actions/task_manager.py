import logging
import atexit
from os import getpid
from heartbeat import HeartBeat
from receiver import Receiver
from params_parser import parse_process_args
from executor.execution_worker import ExecutionWorker

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class taskManager:
    def __init__(self):
        atexit.register(self.handle_shutdown)

        args = parse_process_args()

        self.heartbeat = HeartBeat(process_id=getpid())
        self.heartbeat.daemon = True

        self.receiver = Receiver(cache_prefix=args.execution_group_id, pid=getpid())
        self.executor = ExecutionWorker(args.model_path, args.video_path)

        self.action_map = {
            'merge': self.executor.execute,
            'execute': self.executor.execute
        }

    def handle_shutdown(self):
        self.heartbeat.gracefull_shutdown = True
        logger.info('[TASKMANAGER] Process shutdown successfully')

    def start(self):
        self.heartbeat.start()

        while True:
            job_meta, job_key = self.receiver.get_new_job()
            try:
                self.heartbeat.curr_job_id = job_key.decode('utf-8')
            except Exception as e:
                logger.exception(f'[TASKMANAGER] job key is malformed : {e}')
                continue

            self.heartbeat.set_busy()

            action = self.action_map.get(job_meta['action'])
            if action:
                action(job_meta)
            else:
                logger.warning(f'[TASKMANAGER] action {job_meta["action"]} not defined | ignoring job')

            self.receiver.mark_job_as_done(job_key)
            self.heartbeat.set_free()

        self.heartbeat.join()