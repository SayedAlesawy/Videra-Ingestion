import json
from os import listdir
from os.path import isfile, join
import argparse

def get_ranges(message):  
    encoded_message = []
    i = 0
   
    while (i <= len(message)-1): 
        count = 1
        ch = message[i] 
        j = i 
        while (j < len(message)-1): 
            if (message[j] == message[j+1]): 
                count = count+1
                j = j+1
            else: 
                break
        encoded_message.append([ ch , count])
        i = j+1
    return(encoded_message)



def process_file(file_path):
    res = {}
    with open(file_path,) as f:

        data = json.load(f)
        # print(data)
        keys = list(data.values())
        
        ranges = get_ranges(keys)
        start_idx = 0

        # print(ranges)

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

            # print(tmp, start_idx)
            start_idx += pair[1]
            res.update({tmp: data.get(first)})

        # print("res: ", res)
    return res


def write_file(file_name,res):
    with open('./' + file_name, 'w') as f:
        json.dump(res, f)

def merg_boundry(total_res, res):
    last_idx = list(res.keys())[-1]
    last_value  = res.get(last_idx)

    res_idxs = list(total_res.keys())
    res_last_idx = res_idxs[-1]
    res_last_value = total_res.get(res_last_idx)

    print("in fun: ", last_value, res_last_value)
    print("in fun: ", res)
    print("in fun: ", total_res)
    if(last_value == res_last_value):
        start = res_last_idx.split("-")[0]
        end = last_idx.split("-")[1]
        del res[last_idx]
        del total_res[res_last_idx]
        print("!!!!total: ", total_res)
        print("!!!!!res:", res)
        res.update({last_value: start + "-" + end})
        

    total_res.update(res)
    print("before ret:", total_res)
    return total_res

def parse_process_args():

    parser = argparse.ArgumentParser(description='Read Execution config')
    parser.add_argument('-files-path', required=True)
    parser.add_argument('-no-files', required=True)
    args = parser.parse_args()

    return args.files_path, int(args.no_files)

    

if __name__ == "__main__":
    mypath , number_of_files= parse_process_args()

    # mypath ="../orchestrator/output"
    files_list = [f for f in listdir(mypath) if isfile(join(mypath, f)) and ".json" in f]

    files_list.sort(key= lambda x: int(x.split("-")[1]))


    total_res = {}
    for i, curr_file in enumerate(files_list):

        path = join(mypath, curr_file)
        res = process_file(path)

        if i == 0:
            total_res.update(res)
            continue

        total_res =  merg_boundry(total_res, res)

        total_res.update(res)
        if i == number_of_files:
            break

    write_file("output.json",total_res)