import cv2
import face_recognition
import imutils
import pickle


class ModelExecutor:
    def __init__(self, model_path):
        self.load_model(model_path)

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
                label = "gilfoyle"
                break

        return label

    def load_model(self, model_path):
        """
        user-defined method or we could support standard known formats
        to be loaded by us(h5, ...)
        """
        self.model_instance = pickle.loads(open(model_path, "rb").read())
