import zmq
import json
import cv2


def show_frames(frames):  # utility function for debugging
    for frame in frames:
        cv2.imshow("video", frame)
        cv2.waitKey(5000)


class Receiver:
    def __init__(self, stride):  # should add any additional configs here
        self.stride = stride

    def initializeConnection(self):
        context = zmq.Context()

        socket = context.socket(zmq.REP)
        socket.bind("tcp://*:5555")
        self.socket = socket

    def receiveData(self):
        message = self.socket.recv_json()
        message = json.loads(message)
        return message

    def getFrames(self, info):
        frameIdx = info["frameIndex"]
        batchSize = info["batchSize"]
        cap = cv2.VideoCapture(info["path"])
        totalFrames = cap.get(cv2.CAP_PROP_FRAME_COUNT)

        if batchSize > self.stride:
            # check the frames being read are within boundaries
            # should raise an exception if not
            if frameIdx >= 0 and frameIdx + batchSize <= totalFrames:
                cap.set(cv2.CAP_PROP_POS_FRAMES, frameIdx)
                for i in range(0, batchSize, self.stride):
                    frames = []
                    for i in range(self.stride):
                        ret, frame = cap.read()
                        frames.append(frame)

                    yield frames

    def reply(self):
        self.socket.send(b"World")


if __name__ == "__main__":
    receiver = Receiver(5)
    receiver.initializeConnection()
    while True:
        message = receiver.receiveData()
        print("Received request: %s" % message)
        frame_gen = receiver.getFrames(message)
        while True:
            try:
                frames = next(frame_gen)
                # show_frames(frames)
            except StopIteration:
                print("the end of this patch")
                break

        receiver.reply()
