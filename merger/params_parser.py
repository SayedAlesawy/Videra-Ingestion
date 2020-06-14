import argparse
import logging
from os import path

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()
tag = '[Merger]'


def validate_args(args):
    input_path = args.input_files_path
    if(not path.isdir(input_path)):
        logger.exception(f'{tag} path for files {input_path}\
        passed is not accessible or does not exist')
        return False
    return True


def parse_process_args():
    parser = argparse.ArgumentParser(description='Read Execution config')
    parser.add_argument('-input-files-path', required=True)
    parser.add_argument('-no-files', nargs='?', const=1, type=int, default=-1)
    parser.add_argument('-output-file-path', nargs='?',
                        const=1, type=str, default="output")
    args = parser.parse_args()

    if validate_args(args):
        return args
    else:
        exit(1)
