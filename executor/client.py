import zmq
import json

context = zmq.Context()

print("Connecting to server…")
socket = context.socket(zmq.REQ)
socket.connect("tcp://localhost:5555")

for request in range(10):  # 10 is an arbitrary number
    info = {
        "path": "./mm.mp4",  # just ake any video name in a relative driectory, videos are ignored
        "frameIndex": 60,
        "batchSize": 20
    }
    print("Sending request %s …" % request)
    socket.send_json(json.dumps(info))
    message = socket.recv()
    print("Received reply %s [ %s ]" % (request, message))
