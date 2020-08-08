import datetime
import logging
import mysql.connector
from os import getenv
logging.getLogger().setLevel(logging.INFO)
logger = logging.getLogger()


class DatabaseDriver:
    def __init__(self):
        self.tag = '[DB_DRIVER]'

        self.host = getenv('DB_HOST', 'localhost')
        self.username = getenv('DB_USER', 'root')
        self.passsword = getenv('DB_PASSWORD', 'as')
        self.db_name = getenv('DB_NAME', 'videra_storage')

    def connect(self):
        logger.info(f'{self.tag} establishing a new db connection')
        return mysql.connector.connect(host=self.host,
                                       user=self.username,
                                       password=self.passsword,
                                       database=self.db_name)

    def insert_clips(self, data):
        db_conn = self.connect()
        db_cursor = db_conn.cursor()
        logger.info(f'{self.tag} bulk inserting {len(data)} records')
        sql_statment = "INSERT INTO clips (token, start_time, end_time, tag, updated_at, created_at)\
             VALUES (%s, %s, %s, %s, %s, %s)"

        values = []
        for d in data:
            values.append((d['video_id'], d['start_time'], d['end_time'], d['tag'],
                          str(datetime.datetime.now()), str(datetime.datetime.now())))

        db_cursor.executemany(sql_statment, values)

        db_conn.commit()
