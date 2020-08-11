import logging
import cv2
import json
from os import path
from db.db_driver import DatabaseDriver
logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class Merger():
    def __init__(self, video_path, model_config_path, video_token):
        self.tag = '[MERGER]'
        self.db_driver = DatabaseDriver()
        self.video_path = video_path
        self.video_token = video_token

        self.load_video_meta(video_path)
        self.load_model_config(model_config_path)

    def get_fps(self, video_buffer):
        fps = round(video_buffer.get(cv2.CAP_PROP_FPS))
        logger.info(f"{self.tag} Frames per second: {fps}")
        return fps

    def load_model_config(self, model_config_path):
        with open(model_config_path, "r") as f:
            self.min_period = int(json.load(f).get('min_clip_period')) or 1
            logger.info(f'{self.tag} operating with min_period of {self.min_period}')

    def load_video_meta(self, video_path):
        """
        loads a video at the position of the
        specified frame in the job message
        """
        logger.info(f'{self.tag} loading video meta ....')
        try:
            video_cap = cv2.VideoCapture(video_path)
            self.fps = self.get_fps(video_cap)
            logger.info(f'{self.tag} video cap loaded success')
        except Exception as e:
            logger.exception(f'{self.tag} failed to load video due to : {e}')
            raise

    def find_periods(self, labels_data):
        """
        finds periods with the same label
        """
        min_accepted_period = int(self.min_period * self.fps)
        grouped_data = []
        frames_idxes = sorted(list(labels_data.keys()))
        frames_idxes.append("RESERVED_LABEL")

        current_period_length = 0
        period_start_frame = -1
        current_label = None

        for frame_idx in frames_idxes:
            label = labels_data.get(frame_idx)
            if label == current_label:
                current_period_length += 1
            else:
                if current_period_length > min_accepted_period - 1:
                    grouped_data.append({
                        'video_id': self.video_token,
                        'start_time': int(period_start_frame) / self.fps,
                        'end_time': (int(period_start_frame) + current_period_length) / self.fps,
                        'tag': current_label
                    })
                    current_period_length = 0

                current_label = label
                period_start_frame = frame_idx

        return grouped_data

    def load_data_from_disk(self, labels_file_name):
        logger.info(f'{self.tag} parsing file {labels_file_name}')
        with open(f"./output/{labels_file_name}", "r") as f:
            try:
                return json.load(f)
            except Exception as e:
                logger.error(f'{self.tag} failed to load target file for job {labels_file_name} | error: {e}')
                raise

    def execute(self, job_meta):
        """
        runs the merge execution flow on the current job
        """
        video_file_name = path.basename(self.video_path)
        start_frame_idx = int(job_meta.get('start_idx'))
        frame_end_index = start_frame_idx + int(job_meta.get('frames_count'))
        labels_file_name = f"{video_file_name}-{start_frame_idx}-{frame_end_index}.json"

        labels_data = self.load_data_from_disk(labels_file_name)
        grouped_data = self.find_periods(labels_data)

        self.db_driver.insert_clips(grouped_data)
        logger.info(f'{self.tag} succesffully merged range for file {labels_file_name}')
