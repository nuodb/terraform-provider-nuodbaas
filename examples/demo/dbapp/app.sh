#!/bin/sh

cd "$(dirname "$0")"
pip3 install --quiet --root-user-action=ignore pynuodb
python3 -u app.py | tee out.log
