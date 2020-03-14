import os
import time
import zmq
import logging
import psutil
from threading import Thread

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class HeartBeat(Thread):
    def __init__(self, master_ip='*', master_port=5000, update_frequency=2):
        Thread.__init__(self)

        self.update_frequency = update_frequency

        self.process_id = os.getpid()
        self._total_ram_size = psutil.virtual_memory().total
        self.gracefull_shutdown = False

        self.master_ip = os.getenv('EXECUTION_MANAGER_IP', master_ip)
        self.master_port = os.getenv('EXECUTION_MANAGER_HEARTBEAT_PORT', master_port)

        try:
            self.socket = self.intialize_master_connection()
        except Exception as e:
            logger.error(f'Failed to intialize connection with execution manager due to: {e}')
            raise Exception('CONNECTION_FAILED')

    def intialize_master_connection(self):
        context = zmq.Context()

        socket = context.socket(zmq.PUB)
        socket.bind(f"tcp://{self.master_ip}:{self.master_port}")

        logger.info(f'Connection established successfully with execution manager at port {self.master_port}')
        return socket

    def collect_process_usage_stats(self):
        process_inst = psutil.Process(self.process_id)

        ram_usage = process_inst.memory_info()[0]
        ram_usage = ram_usage / self._total_ram_size

        try:
            gpu_usage = self._get_gpu_usage_by_process()
        except Exception:
            gpu_usage = '0'

        cpu_percent = process_inst.cpu_percent()
        # return f"pid:{self.process_id}|cpu:{cpu_percent}%|ram:{ram_usage}%|gpu:{gpu_usage}%"
        return {"pid": self.process_id, "cpu": cpu_percent, "ram": ram_usage, "gpu": gpu_usage}

    def send_heartbeat(self):
        usage_stats = self.collect_process_usage_stats()

        logger.info(f'sending usage_stats to execution manager: {usage_stats}')
        self.socket.send_json(usage_stats)

    def _get_gpu_usage_by_process(self):
        """
        THIS FUNCTION/CODE WAS ADOPTED FROM A BROKEN PYTHON PKG
        THE PKG WAS INTENTED TO CHECK GPU STATS
        LINK: https://github.com/FlyHighest/gpuinfo
        """
        ns = os.popen('nvidia-smi')
        lines_ns = ns.readlines()
        total_gpu_memory = 0
        for line in lines_ns:
            if('%' in line):
                total_gpu_memory = int(line.split('MiB')[-2].replace('/', ''))
            if(str(self.process_id) in line and total_gpu_memory > 0):
                return int(line.split(' ')[-2].replace('MiB', '')) / total_gpu_memory

        return '0'  # if no entry for us then this process not using gpu

    def run(self):
        logger.info(f'Heartbeat thread started on process with id-{self.process_id}')

        while time.sleep(self.update_frequency) or not self.gracefull_shutdown:
            try:
                self.send_heartbeat()
            except Exception as e:
                logger.error(f'Failed to send heartbeat to execution manager on port {self.master_port} due to: {e}')

        # self.socket.disconnect(f"tcp://{self.master_ip}:{self.master_port}")
