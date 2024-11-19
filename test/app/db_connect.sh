#!/usr/bin/env sh

set -e

# Setup Python virtual environment
python3 -m venv .venv
. .venv/bin/activate
pip3 install -r requirements.txt

python3 db_connect.py
