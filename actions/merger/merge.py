import logging
import cv2
import json
from os import path
from db.db_driver import DatabaseDriver
logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class Merger():
    def __init__(self, video_path):
        self.tag = '[MERGER]'
        self.db_driver = DatabaseDriver()
        self.video_path = video_path
        self.load_video_meta()

    def get_fps(self, video_buffer):
        fps = round(video_buffer.get(cv2.CAP_PROP_FPS))
        print(f"Frames per second: {fps}")
        return fps

    def load_video_meta(self, video_path):
        """
        loads a video at the position of the
        specified frame in the job message
        """
        logger.info(f'{self.tag} loading video ....')
        try:
            video_cap = cv2.VideoCapture(video_path)
            self.fps = self.get_fps(video_cap)
            logger.info(f'{self.tag} video cap loaded success')
        except Exception as e:
            logger.exception(f'failed to load video due to : {e}')

    def find_periods(self, labels_data, labels_file_name, base_frame_idx):
        """
        finds periods with the same label
        """
        min_accepted_period = int(5 * self.fps)
        grouped_data = []
        frames_idxes = sorted(list(labels_data.keys()))

        current_period_length = 0
        period_start_frame = -1
        current_label = None
        for frame_idx in frames_idxes:
            label = labels_data[frame_idx]
            if label == current_label:
                current_period_length += 1
            else:
                if current_period_length > min_accepted_period - 1:
                    grouped_data.append({
                        'video_id': labels_file_name,
                        'start_time': base_frame_idx + period_start_frame,
                        'end_time': base_frame_idx + period_start_frame + current_period_length,
                        'tag': current_label
                    })

                current_label = label
                period_start_frame = frame_idx

        return grouped_data

    def execute(self, job_meta):
        """
        runs the merge execution flow on the current job
        """
        video_file_name = path.basename(self.video_path)
        start_frame_idx = int(job_meta.get('start_idx'))
        frame_end_index = start_frame_idx + int(job_meta.get('frames_count'))
        labels_file_name = f"{video_file_name}-{start_frame_idx}-{frame_end_index}.json"

        with open(f"./output/{labels_file_name}", "r") as f:
            labels_data = json.loads(labels_file_name, f)
            grouped_data = self.find_periods(labels_data, labels_file_name, start_frame_idx)
            self.db_driver.insert_clips(grouped_data)
