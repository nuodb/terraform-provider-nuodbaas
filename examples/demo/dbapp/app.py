"""
Connect to a NuoDB database and run a basic workload.

The following environment variables must be set:
 - DB_NAME: The name of the database
 - DB_USER: A database user
 - DB_PASSWORD: The database user's password
 - DB_HOST: Hostname and port for the database server
 - CA_CERT: The PEM-encoded certificate to verify the database server
"""

import os
import time

import pynuodb


class SimpleApp(object):
    def __init__(self, db: str, user: str, password: str, host: str, cert: str):
        self.db = db
        self.user = user
        self.password = password
        self.host = host

        self.cert_file = "/tmp/cert.pem"
        with open(self.cert_file, mode="w") as f:
            f.write(cert)
            f.flush()

        self.connection = None
        self.cursor = None

    def connect(self):
        self.close()
        self.connection = pynuodb.connect(
            database=self.db,
            user=self.user,
            password=self.password,
            host=self.host,
            options={"trustStore": self.cert_file},
        )
        self.cursor = self.connection.cursor()

    def close(self):
        if self.cursor:
            self.cursor.close()
            self.cursor = None
        if self.connection:
            self.connection.close()
            self.connection = None

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_value, traceback):
        self.close()

    def create_schema(self):
        self.cursor.execute("drop table if exists testdata")
        self.cursor.execute("create table testdata ( time integer )")
        self.connection.commit()

    def run(self):
        print("Connecting to database and creating table...")
        self.connect()
        self.create_schema()

        i = 0
        while True:
            i += 1
            if i % 10 == 0:
                print("Reconnecting to database...")
                self.connect()
            try:
                current_time = int(time.time())
                print("Writing seconds since epoch {}".format(current_time))
                self.cursor.execute("insert into testdata values (?)", (current_time,))
                self.connection.commit()
            except Exception as e:
                print("Unable to write data: " + str(e))
            # Insert delay before next iteration
            time.sleep(1)


if __name__ == "__main__":
    db = os.getenv("DB_NAME")
    user = os.getenv("DB_USER")
    password = os.getenv("DB_PASSWORD")
    host = os.getenv("DB_HOST")
    cert = os.getenv("CA_CERT")
    with SimpleApp(db, user, password, host, cert) as app:
        app.run()
