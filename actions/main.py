import sys
import logging
from os import getpid
from task_manager import taskManager

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

stream = logging.StreamHandler(sys.stdout)
stream.setLevel(logging.INFO)
logger.addHandler(stream)

fh = logging.FileHandler(f'./logs/executor/logs-{getpid()}.log')
fh.setLevel(logging.INFO)
logger.addHandler(fh)

if __name__ == "__main__":

    logger.info(f'[EXEC]  Model executor process started with id-{getpid()}')

    tm = taskManager()
    tm.start()
