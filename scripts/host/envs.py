from common import *
import os

CLOUDLAB_USER = os.getenv("node_username")
if CLOUDLAB_USER is None:
    exit("node_username is not set in environment. source env.sh first.")
KEYFILE = os.getenv("private_key")
if KEYFILE is None:
    exit("private_key is not set in environment. source env.sh first.")

HOST_SERVERS = [f'{CLOUDLAB_USER}@{s}' for s in SERVERS]
