import argparse
import logging
from os import path

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()

tag = '[executor-parser]'


def validate_args(args):
    files = {'model_path': args.model_path, 'video_path': args.video_path, 'config_path': args.model_config_path}
    for argname in files:
        fpath = files[argname]
        if not path.isfile(fpath):
            logger.exception(f'{tag} File Path for arg {argname} passed is not accessible {fpath}') # noqa
            return False
    return True


def parse_process_args():
    parser = argparse.ArgumentParser(description='Read Execution config')
    parser.add_argument('-video-path', help='target video path on disk', required=True)
    parser.add_argument('-model-path', help='model.h5 file path to load model weights', required=True)
    parser.add_argument('-model-config-path',
                        help='contains info for how the model input is expected',
                        required=True)
    parser.add_argument('-execution-group-id', help='id asssigned to the execution manager, used to scope redis queues')
    parser.add_argument('-code-path', help='path to the code file provided by the user')
    parser.add_argument('-video-token', help='video identitfying token in database')
    args = parser.parse_args()

    if validate_args(args):
        return args
    else:
        exit(1)
