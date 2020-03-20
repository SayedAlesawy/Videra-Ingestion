import logging

logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class ExecutionWorker():
    def __init__(self, model_path):
        self.model_instance = self.load_model(model_path)
        self.labels = []

    def load_model(self, model_path):
        """
        user-defined method or we could support standard known formats
        to be loaded by us(h5, ...)
        """
        logger.info('loading model ....')
        logger.info('model loaded success')
        pass

    def execute_packets(self, packets):
        logger.info(f'recieved {packets} packets to execute')
        processed_packets_count = 0
        try:
            for packet in packets:
                processed_packets_count += 1
                self.execute_packet(packet)
        except StopIteration:
            logger.info(f'succesfully processed {processed_packets_count}')

        return self.labels

    def execute_packet(self, data_packet):
        try:
            lbl = self.run_model(data_packet)
            self.labels.append(lbl)
        except Exception as e:  # the model may fail for many reasons, we need to handle failures in labeling
            logger.error(f'model failed to label data packet due to: {e}')
            logger.debug(data_packet)
            self.labels.append(None)

    def run_model(self, data_packet):
        """
        user-defined method
        model is avaiable through self
        expected output: label(string, numeric)
        """
        return 0
