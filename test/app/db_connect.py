"""
Try connecting to a NuoDB database.

It assumes the following environment variables are set:
 - DB_NAME: The name of the database
 - DB_USER: A database user
 - DB_PASSWORD: The database user's password
 - PEER_ADDRESS: Hostname and port for the database server
 - CA_CERT: The PEM-encoded certificate to verify the database server
"""

import os
import tempfile

import pynuodb


def connect(db: str, user: str, password: str, url: str, cert: str):
    print("Creating cert file.")
    with tempfile.NamedTemporaryFile(mode="w", prefix="cert", suffix=".pem") as cert_file:
        cert_file.write(cert)
        cert_file.flush()

        print("Connecting to the database.")
        connection = pynuodb.connect(database=db,
                                     user=user,
                                     password=password,
                                     host=url,
                                     options={"trustStore": cert_file.name})
        
        print("Connected, testing a basic query.")
        cursor = connection.cursor()
        try:
            cursor.execute("SELECT 1 FROM DUAL")
            rows = cursor.fetchall()
            assert len(rows) == 1
            assert rows[0] == (1,)
        finally:
            cursor.close()
            connection.close()
        print("Query succeeded")


if __name__ == '__main__':
    db = os.getenv("DB_NAME")
    user = os.getenv("DB_USER")
    password = os.getenv("DB_PASSWORD")
    url = os.getenv("PEER_ADDRESS")
    cert = os.getenv("CA_CERT")
    connect(db, user, password, url, cert)
