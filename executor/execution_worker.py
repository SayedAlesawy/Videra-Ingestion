import logging
import cv2
import face_recognition
import imutils
import pickle
import json
from os import path
logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class ExecutionWorker():
    def __init__(self, model_path, video_path):
        self.tag = '[EXECUTION_WORKER]'
        self.labels = {}
        self.model_path = model_path
        self.video_path = video_path

        self.load_model(model_path)
        self.stride = 1

    def load_model(self, model_path):
        """
        user-defined method or we could support standard known formats
        to be loaded by us(h5, ...)
        """
        logger.info(f'{self.tag} loading model ....')
        try:
            self.model_instance = pickle.loads(open(model_path, "rb").read())
            logger.info(f'{self.tag} model loaded success')
        except Exception as e:
            logger.exception(f'failed to load model due to : {e}')

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
            logger.exception(f'failed to load video due to : {e}')

    def write_labels(self):
        """
        write the output of execution as json formatted
        {frame_id: label}
        to disk
        """
        video_file_name = path.basename(self.video_path)
        labels_file_name = f"{video_file_name}-{self.start_frame_idx}-{self.frame_end_index}.json"
        logger.info(f"{self.tag} writing output to {labels_file_name}")
        with open(f"./output/{labels_file_name}", "w") as f:
            json.dump(self.labels, f)

        self.labels.clear()

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
            lbl = self.run_model(data_packet)
            self.labels[self.start_frame_idx + index] = lbl
        except Exception as e:  # the model may fail for many reasons, we need to handle failures in labeling
            logger.error(f'model failed to label data packet due to: {e}')
            logger.debug(data_packet)
            self.labels[self.start_frame_idx + index] = None

    def run_model(self, data_packet):
        """
        user-defined method
        model is avaiable through self
        expected output: label(string || numeric)
        """
        frame = data_packet[0]
        rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        rgb = imutils.resize(frame, width=750)

        # detect the (x, y)-coordinates of the bounding boxes
        # corresponding to each face in the input frame, then compute
        # the facial embeddings for each face
        boxes = face_recognition.face_locations(rgb, model='hog')
        encodings = face_recognition.face_encodings(rgb, boxes)

        label = False
        # loop over the facial embeddings
        for encoding in encodings:
            # attempt to match each face in the input image to our known
            # encodings
            matches = face_recognition.compare_faces(self.model_instance["encodings"], encoding)
            label = True in matches
            if(label):
                break

        return label
