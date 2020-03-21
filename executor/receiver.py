import zmq
import json
import cv2
import logging

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


def show_frames(frames):  # utility function for debugging
    for frame in frames:
        cv2.imshow("video", frame)
        cv2.waitKey(5000)


class Receiver:
    def __init__(self, stride):  # should add any additional configs here
        self.stride = stride
        self.metadata = {}
        self.initializeConnection()

    def initializeConnection(self, port=5555, host='*'):
        context = zmq.Context()

        socket = context.socket(zmq.REP)
        socket.bind(f"tcp://{host}:{port}")
        logger.info(f'Established connection with master on tcp://{host}:{port}')
        self.socket = socket

    def model_path(self):
        return self.metadata.get('path')

    def get_batch_metadata(self):
        message = self.socket.recv_json()
        if type(message) is dict:
            metadata = message
        else:
            metadata = json.loads(message)

        logger.info("Received Job Meta: %s" % metadata)
        self.metadata = metadata

    def validate_metadata(self):
        info = self.metadata
        frameIdx = info["frameIndex"]
        batchSize = info["batchSize"]

        cap = cv2.VideoCapture(info["path"])
        totalFrames = cap.get(cv2.CAP_PROP_FRAME_COUNT)
        logger.info('validating job metadata')

        if batchSize > self.stride and frameIdx >= 0 and frameIdx + batchSize <= totalFrames:
            cap.set(cv2.CAP_PROP_POS_FRAMES, frameIdx)
            logger.info('metadata valid, parsing batch')
            return cap
        else:
            logger.info('batch discared due to inconsistent batch info')
            return False

    def generate_data(self):
        self.get_batch_metadata()
        cap = self.validate_metadata()
        if(cap):
            for i in range(0, self.metadata['batchSize'], self.stride):
                frames = []
                for i in range(self.stride):
                    ret, frame = cap.read()
                    frames.append(frame)

                yield frames

        self.reply()

    def reply(self):
        self.socket.send(b"World")
