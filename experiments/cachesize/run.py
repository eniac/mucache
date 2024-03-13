#!/usr/bin/env python3

import signal
import time
from experiments.helper import *

APP = "hotel"
set_app(APP)

def start_proxy():
    run_shell("cd proxy && cargo build --release")
    frontend_ip = get_ip("frontend")
    p = run_in_bg(
        f"cargo run --release hotel --frontend {frontend_ip}",
        "proxy")
    time.sleep(5)
    return p

def populate():
    args = ""
    for service in ["frontend", "user"]:
        ip = get_ip(service)
        args += f" --{service} {ip}"
    args += f" --info_size 5000"
    run_shell("python3 experiments/hotel/populate.py" + args)


def run_once(mem: int, cm: str = "true"):
    clean2(mem)
    deploy(cm=cm)
    populate()
    p = start_proxy()
    # warm up
    run_shell(compose_oha_proxy(req=1000, duration=120))
    time.sleep(2)
    old = get_hit_miss_redis()
    # run
    res = run_shell(compose_oha_proxy(req=1000, duration=60))
    res = parse_res(res)
    new = get_hit_miss_redis()
    os.kill(p.pid, signal.SIGINT)
    p.terminate()
    p.wait()
    res["hit_rate"] = compute_hit_rate_redis(old, new)
    return res


def main():
    # 2000 reqs, 120s takes around 128MB in home_timeline
    mems = [16, 32, 64, 128, 256, 512, 1024]

    ours = {}
    for mem in mems:
        ours[mem] = run_once(mem)
    with open(f"{APP}-md.json", "w") as f:
        json.dump(ours, f, indent=2)

    uppers = {}
    for mem in mems:
        uppers[mem] = run_once(mem, cm="upper")
    with open(f"{APP}-md-upper.json", "w") as f:
        json.dump(uppers, f, indent=2)
    clean2(mem)

    print(ours)
    print(uppers)


if __name__ == "__main__":
    main()
