import time
import zmq
import json
import cv2

#later will be sent through process args, for now kept as a constant
NUMBER_OF_FRAMES = 3


#might change this pattern
def initializeConnection():
    context = zmq.Context()

    socket = context.socket(zmq.REP)
    socket.bind("tcp://*:5555")
    return socket


def receiveData(socket):
    message = socket.recv_json()
    message = json.loads(message)
    return message

def getFrames(info):
    frameIdx = info["frameIndex"]
    cap = cv2.VideoCapture(info["path"])
    totalFrames = cap.get(cv2.CAP_PROP_FRAME_COUNT)
    
    #check the frames being read are within boundaries
    #should raise an exception if not
    if frameIdx >= 0 and frameIdx + NUMBER_OF_FRAMES <= totalFrames:
        cap.set(cv2.CAP_PROP_POS_FRAMES,frameIdx)
        ret, frame = cap.read()
        for i in range(NUMBER_OF_FRAMES):
            yield frame
            ret, frame = cap.read()

def reply():
    socket.send(b"World")


if __name__ == "__main__":
    socket = initializeConnection()

    while True:
        message = recieveData(socket)
        print("Received request: %s" % message)
        
        frame_gen = getFrames(message)
        while True:
            try:
                cv2.imshow("video", next(frame_gen))
                cv2.waitKey(5000)
                #processing the frames could be here
            except StopIteration:
                print("the end of this patch")
                break

        reply()