import redis
import time
import json
import logging
import os
logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class FlowManager:
    def __init__(self, cache_prefix, pid, redis_host=os.getenv('REDIS_HOST'), redis_port=6379, update_frequency=1):
        self.tag = '[RECIEVER]'
        prefix = f'{cache_prefix}:ingestion'
        self.todo_list = f'{prefix}:todo'
        self.inprogress_list = f'{prefix}:in-progress'
        self.done_list = f'{prefix}:done'
        self.jobs = f'{prefix}:jobs'
        self.active_jobs = f'{prefix}:active_jobs'
        self.handshake = f'{prefix}:ready'

        self.is_ready = False
        self.update_frequency = update_frequency
        self.pid = str(pid)

        self.initializeConnection(redis_host, redis_port)

    def initializeConnection(self, redis_host, redis_port):
        logger.info(f'{self.tag} intializing connection with redis on host: {redis_host}, port: {redis_port}')
        try:
            self.redis_instance = redis.Redis(host=redis_host, port=redis_port)
            logger.info(f'{self.tag} redis instance established')
        except Exception as err:
            logger.exception(f'{self.tag} failed to intialize connection with redis due to: {err}')
            raise

    def request_handshake(self):
        logger.info(f'{self.tag} requesting handshak on key set {self.handshake} for pid {self.pid}')
        self.redis_instance.hset(self.handshake, self.pid, 'False')
        logger.info(f'{self.tag} requested handshak on key set {self.handshake} for pid {self.pid} successfully')

    def check_handshake(self):
        logger.info(f'{self.tag} polling on handshake')
        while time.sleep(self.update_frequency):
            current_state = self.redis_instance.hget(self.handshake, self.pid).decode('utf-8')
            logger.info(f'{self.tag} handshake state: {current_state}')

            if current_state == 'True':
                self.request_handshake()
                return

    def get_new_job(self):
        """
        checks todo queue for new jobs
        moves found job from todo queue to in-progress and
        sets the current active job to the found job
        fetches the job info from the jobs mapping set
        returns both the job metainfo and the job key
        """
        self.check_handshake()

        while True:
            job_meta_string = None
            job_md5_hash = None
            try:
                logger.info(f'{self.tag} checking todo list for tasks ...')
                job_md5_hash = self.redis_instance.brpoplpush(self.todo_list, self.inprogress_list, timeout=0)
                self.redis_instance.hset(self.active_jobs, self.pid, job_md5_hash)
                job_meta_string = self.redis_instance.hget(self.jobs, job_md5_hash)
                logger.info(f'{self.tag} got new job {job_md5_hash}, parsing job ...')
            except Exception as e:
                logger.exception(f'{self.tag} failed to get job meta from redis due to: {e}')
                continue

            try:
                job_meta = json.loads(job_meta_string)
                logger.info(f'{self.tag} job parsed successfully')
                return job_meta, job_md5_hash
            except Exception as err:
                logger.exception(f'{self.tag} failed to process job meta due to: {err}')

    def reject_job(self, job_key):
        """
        sends the job back to todo list
        """
        atomic_pipeline = self.redis_instance.pipeline()

        atomic_pipeline.lrem(self.inprogress_list, 0, job_key)
        atomic_pipeline.rpush(self.todo_list, job_key)
        atomic_pipeline.hdel(self.active_jobs, self.pid)
        atomic_pipeline.execute()
