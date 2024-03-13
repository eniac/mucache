import time

from fabric import Connection, ThreadingGroup
from common import *
from envs import *

SETUP_PATH = "./mucache/scripts/setup"


def master_conn():
    return Connection(HOST_SERVERS[0], connect_kwargs={
        'key_filename': KEYFILE,
    })


def worker_conn(idx: int):
    return Connection(HOST_SERVERS[idx], connect_kwargs={
        'key_filename': KEYFILE,
    })


def setup_master():
    conn = master_conn()
    res = conn.run(os.path.join(SETUP_PATH, "master.sh"))
    print(res)


def setup_workers():
    workers = HOST_SERVERS[1:]
    group = ThreadingGroup(*workers, connect_kwargs={
        'key_filename': KEYFILE,
    })
    res = group.run(os.path.join("./", "worker.sh"))
    print(res)


def start_dapr():
    conn = master_conn()
    res = conn.run(os.path.join(SETUP_PATH, "start_dapr.sh") + f" {len(HOST_SERVERS) - 1}")
    print(res)


def setup():
    setup_master()
    setup_workers()
    time.sleep(3)
    start_dapr()


def main():
    setup()


if __name__ == '__main__':
    main()
