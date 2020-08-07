import logging
import atexit
from os import getpid
from heartbeat import HeartBeat
from flow_manager import FlowManager
from params_parser import parse_process_args
from executor.execution_worker import ExecutionWorker
from merger.merge import Merger

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class taskManager:
    def __init__(self):
        args = parse_process_args()

        self.heartbeat = HeartBeat(process_id=getpid())
        self.heartbeat.daemon = True
        atexit.register(self.handle_shutdown)

        self.flow_manager = FlowManager(cache_prefix=args.execution_group_id, pid=getpid())
        self.executor = ExecutionWorker(args.model_path, args.video_path, args.model_config_path, './actions/executor')
        self.merger = Merger(args.video_path, args.model_config_path)

        self.action_map = {
            'merge': self.merger.execute,
            'execute': self.executor.execute
        }

    def handle_shutdown(self):
        self.heartbeat.gracefull_shutdown = True
        logger.info('[TASKMANAGER] Process shutdown successfully')

    def start(self):
        self.heartbeat.start()

        while True:
            try:
                job_meta, job_key = self.flow_manager.get_new_job()
                self.heartbeat.curr_job_id = job_key.decode('utf-8')
            except Exception as e:
                logger.exception(f'[TASKMANAGER] job key is malformed : {e}')
                self.flow_manager.reject_job(job_key)
                continue

            self.heartbeat.set_busy()

            action = self.action_map.get(job_meta['action'])
            if action:
                try:
                    action(job_meta)

                    self.flow_manager.mark_job_as_done(job_key)
                except Exception as e:
                    logger.exception(f'[TASKMANAGER] failed to execute action on job {job_meta.get("jid")} | {e}')
                    self.flow_manager.reject_job(job_key)
            else:
                logger.warning(f'[TASKMANAGER] action {job_meta.get("action")} not defined | ignoring job')

            self.heartbeat.set_free()

        self.heartbeat.join()
