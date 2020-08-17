import logging
import cv2
import json
import sys
from os import path
from collections import defaultdict
logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class ExecutionWorker():
    def __init__(self, model_path, video_path, model_config_path, code_path):
        self.tag = '[EXECUTION_WORKER]'
        self.labels = defaultdict(list)
        self.model_path = model_path
        self.video_path = video_path

        self.load_model_config(model_config_path)

        sys.path.append(code_path)
        from code_file import ModelExecutor # noqa
        self.model_exec = ModelExecutor(model_path)

    def load_model_config(self, model_config_path):
        with open(model_config_path, "r") as f:
            self.stride = int(json.load(f).get('stride')) or 1
            logger.info(f'{self.tag} operating with a stride of {self.stride}')

    def load_video(self, video_path):
        """
        loads a video at the position of the
        specified frame in the job message
        """
        logger.info(f'{self.tag} loading video ....')
        try:
            self.video_cap = cv2.VideoCapture(video_path)
            self.video_cap .set(cv2.CAP_PROP_POS_FRAMES, self.start_frame_idx)
            logger.info(f'{self.tag} video cap loaded success')
        except Exception as e:
            logger.exception(f'{self.tag} failed to load video from {video_path} due to : {e}')
            raise

    def write_labels(self):
        """
        write the output of execution as json formatted
        {frame_id: label}
        to disk
        """
        video_file_name = path.basename(self.video_path)
        labels_file_name = f"{video_file_name}-{self.start_frame_idx}-{self.frame_end_index}.json"
        logger.info(f"{self.tag} writing output to {labels_file_name}")

        try:
            with open(f"./output/{labels_file_name}", "w") as f:
                json.dump(self.labels, f)
        except Exception as e:
            logger.exception(f'{self.tag} failed to write labels output for video {labels_file_name} due to: {e}')
            raise

        self.labels.clear()

    def check_if_written(self):
        video_file_name = path.basename(self.video_path)
        labels_file_name = f"{video_file_name}-{self.start_frame_idx}-{self.frame_end_index}.json"
        return path.isfile(f"./output/{labels_file_name}")

    def execute(self, job_meta):
        """
        executes given job by calling appropiate methods in
        order
        """
        self.start_frame_idx = int(job_meta.get('start_idx'))
        self.frame_end_index = self.start_frame_idx + int(job_meta.get('frames_count'))
        self.load_video(self.video_path)
        packets = self.load_packets()
        self.execute_packets(packets)
        self.write_labels()

    def execute_packets(self, packets):
        """
        executes packets produced by load_packets concurrently
        """
        logger.info(f'{self.tag} recieved packets to execute (consumer-producer)')
        processed_packets_count = 0
        for packet in packets:
            self.execute_packet(packet, processed_packets_count)
            processed_packets_count += 1

        logger.info(f'{self.tag} succesfully processed {processed_packets_count}')
        return self.labels

    def load_packets(self):
        """
        loads frames from disk and
        groups them into arrays of size stride
        each packets contains x frames as specified
        by model config
        """
        packets = []
        logger.info(f'{self.tag} loading packets from vidoe source (consumer-producer)')
        for i in range(self.start_frame_idx, self.frame_end_index, self.stride):
            frames = []
            for j in range(self.stride):
                ret, frame = self.video_cap.read()
                frames.append(frame)
            packets.append(frames)
            yield frames

    def execute_packet(self, data_packet, index):
        try:
            lbl = self.model_exec.run_model(data_packet)
            if lbl:
                self.labels[lbl].append(self.start_frame_idx + index)
        except Exception as e:  # the model may fail for many reasons, we need to handle failures in labeling
            logger.error(f'model failed to label data packet due to: {e}')
            logger.debug(data_packet)
            self.labels['failures'].append(self.start_frame_idx + index)
