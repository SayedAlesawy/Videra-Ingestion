import mysql.connector
from os import getenv


class DatabaseDriver:
    def __init__(self):
        self.host = getenv('DB_HOST')
        self.username = getenv('DB_USER')
        self.passsword = getenv('DB_PASSWORD')
        self.db_name = getenv('DB_NAME')

    def connect(self):
        return mysql.connector.connect(host=self.host,
                                       user=self.username,
                                       password=self.passsword,
                                       database=self.db_name)

    def init_schema(self):
        """
        example
        clip:
            video_id: mm1
            start_time: 00:11:23
            end_time: 00:11:33
            tag: anything
        """
        db_conn = self.connect()
        db_cursor = db_conn.cursor()
        db_cursor.execute(
            "CREATE TABLE clips (id INT AUTO_INCREMENT PRIMARY KEY, \
                video_id VARCHAR(255), start_time VARCHAR(255), end_time, tag VARCHAR(255))"
        )

    def insert_clips(self, data):
        db_conn = self.connect()
        db_cursor = db_conn.cursor()

        sql_statment = "INSERT INTO clips (video_id, start_time, end_time, tag) VALUES (%s, %s, %s, %s)"
        values = []
        for d in data:
            values.append(d['video_id'], d['start_time'], d['end_time'], d['tag'])

        db_cursor.executemany(sql_statment, values)

        db_conn.commit()
