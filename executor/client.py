import zmq
import json

context = zmq.Context()

print("Connecting to server…")
socket = context.socket(zmq.REQ)
socket.connect("tcp://localhost:5555")

for request in range(10):  # 10 is an arbitrary number
    info = {
        "path": "/home/yasmeen/Desktop/video.mkv",
        "frameIndex": 34000,
        "batchSize": 20
    }
    print("Sending request %s …" % request)
    socket.send_json(json.dumps(info))
    message = socket.recv()
    print("Received reply %s [ %s ]" % (request, message))
