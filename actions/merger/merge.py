import json
import logging
import sys
from os import listdir, stat, mkdir, getpid
from os.path import isfile, join, dirname
from params_parser import parse_process_args

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()
tag = '[Merger]'


def get_ranges(run_count):
    run = []
    i = 0

    while (i <= len(run_count)-1):
        count = 1
        ch = run_count[i]
        j = i
        while (j < len(run_count)-1):
            if (run_count[j] == run_count[j+1]):
                count = count+1
                j = j+1
            else:
                break
        run.append([ch, count])
        i = j+1
    return(run)


def process_file(file_path):
    res = {}

    with open(file_path,) as f:

        data = json.load(f)
        keys = list(data.values())

        ranges = get_ranges(keys)
        start_idx = 0

        for i, pair in enumerate(ranges):

            first = start_idx
            second = pair[1] + start_idx

            if(i == 0):
                second -= 1

            data_keys = list(data.keys())

            if(second == len(data_keys)):
                second -= 1

            first = data_keys[first]
            second = data_keys[second]

            tmp = str(first) + "-" + str(second)

            start_idx += pair[1]
            res.update({tmp: data.get(first)})

    return res


def write_file(file_path, res):
    directory = dirname(file_path)
    try:
        stat(directory)
    except(Exception):
        mkdir(directory)

    with open('./' + file_path, 'w') as f:
        json.dump(res, f)


def merg_boundry(total_res, res):
    last_idx = list(res.keys())[-1]
    last_value = res.get(last_idx)

    res_idxs = list(total_res.keys())
    res_last_idx = res_idxs[-1]
    res_last_value = total_res.get(res_last_idx)

    if(last_value == res_last_value):
        start = res_last_idx.split("-")[0]
        end = last_idx.split("-")[1]
        del res[last_idx]
        del total_res[res_last_idx]

        res.update({last_value: start + "-" + end})

    total_res.update(res)
    return total_res


if __name__ == "__main__":

    args = parse_process_args()
    mypath = args.input_files_path
    number_of_files = args.no_files
    output_path = args.output_file_path

    stream = logging.StreamHandler(sys.stdout)
    stream.setLevel(logging.INFO)
    logger.addHandler(stream)

    fh = logging.FileHandler(f'./logs/merger/logs-{getpid()}.log')
    fh.setLevel(logging.INFO)
    logger.addHandler(fh)

    files_list = [f for f in listdir(mypath)
                  if isfile(join(mypath, f)) and ".json" in f]

    files_list.sort(key=lambda x: int(x.split("-")[1]))

    total_res = {}
    actual_number_of_files = len(files_list)
    for i, curr_file in enumerate(files_list):

        path = join(mypath, curr_file)
        res = process_file(path)

        if i == 0:
            total_res.update(res)
            continue

        total_res = merg_boundry(total_res, res)

        total_res.update(res)
        logger.info(f'{tag} finished processing {i+1} files out of\
                    {actual_number_of_files}')
        if i == number_of_files:
            break

    output_path = join(output_path, curr_file.split(".")[0] + ".json")
    write_file(output_path, total_res)