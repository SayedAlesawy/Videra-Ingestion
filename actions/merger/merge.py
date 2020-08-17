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
            model_config = json.load(f)
            self.min_period = float(model_config.get('min_clip_period')) or 1  # in secs
            self.miss_tolerance = int(model_config.get('frame_miss_tolerance')) or 5  # in frames
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

    def check_period(self, min_accepted_period, group_end, group_start, grouped_data, label):
        if min_accepted_period <= group_end - group_start + 1:
            grouped_data.append({
                'video_id': self.video_token,
                'start_time': group_start / self.fps,
                'end_time': max(group_end / self.fps, group_start / self.fps + 1),
                'tag': label
            })

    def find_periods(self, labels_data):
        """
        finds periods with the same label
        """
        min_accepted_period = int(self.min_period * self.fps)
        grouped_data = []
        labels_set = list(labels_data.keys())

        for label in labels_set:
            frames_with_label = labels_data.get(label)
            if(frames_with_label and len(frames_with_label)):
                frames_with_label = sorted(frames_with_label)
                group_start = int(frames_with_label[0])
                group_end = group_start

                for frame in frames_with_label[1:]:
                    if abs(frame - group_end) <= self.miss_tolerance:
                        # add new member to group
                        group_end = int(frame)
                    else:
                        # group has ended, let's check it passes min perid check
                        self.check_period(min_accepted_period, group_end, group_start, grouped_data, label)
                        # start a new period
                        group_start = int(frame)
                        group_end = group_start

            self.check_period(min_accepted_period, group_end, group_start, grouped_data, label)
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
