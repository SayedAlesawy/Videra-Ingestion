import zmq
import json
import cv2


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
        cap = cv2.VideoCapture(info["path"])
        totalFrames = cap.get(cv2.CAP_PROP_FRAME_COUNT)

        # check the frames being read are within boundaries
        # should raise an exception if not
        if frameIdx >= 0 and frameIdx + self.stride <= totalFrames:
            cap.set(cv2.CAP_PROP_POS_FRAMES, frameIdx)
            ret, frame = cap.read()
            for i in range(self.stride):
                yield frame
                ret, frame = cap.read()

    def reply(self):
        self.socket.send(b"World")


if __name__ == "__main__":
    receiver = Receiver(3)
    receiver.initializeConnection()
    while True:
        message = receiver.receiveData()
        print("Received request: %s" % message)
        frame_gen = receiver.getFrames(message)
        while True:
            try:
                cv2.imshow("video", next(frame_gen))
                cv2.waitKey(5000)
                # processing the frames could be here
            except StopIteration:
                print("the end of this patch")
                break

        receiver.reply()
