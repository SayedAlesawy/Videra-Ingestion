import argparse
import os
import logging

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

tag = '[executor-parser]'


def validate_args(args):
    if not os.path.isfile(args.model_path) or not os.path.isfile(args.video_path) or not os.path.isfile(
            args.model_config_path):
        logger.exception(f'{tag} File Paths passed are not accessible {args.model_path} | {args.video_path} | {args.model_config_path}') # noqa
        return False
    return True


def parse_process_args():
    parser = argparse.ArgumentParser(description='Read Execution config')
    parser.add_argument('-video-path', help='target video path on disk', required=True)
    parser.add_argument('-model-path', help='model.h5 file path to load model weights', required=True)
    parser.add_argument('-model-config-path',
                        help='contains info for how the model input is expected',
                        required=True)
    args = parser.parse_args()

    if validate_args(args):
        return args
    else:
        exit()
