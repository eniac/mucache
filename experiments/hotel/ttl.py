#!/usr/bin/env python3
import signal
import time
from experiments.helper import *
from collections import defaultdict
from pprint import pprint

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
    run_shell("python3 experiments/hotel/populate.py" + args)


def run_once(req: int, cm: str, ttl=None):
    clean2(mem="20")
    deploy(cm=cm, ttl=ttl)
    populate()
    p = start_proxy()
    top_p, top_q = top_process()
    res = run_shell(compose_oha_proxy(req=req, duration=120))
    res = parse_res(res)
    os.kill(p.pid, signal.SIGINT)
    p.terminate()
    p.wait()
    if cm in ["true", "upper"]:
        res["hit_rate"] = get_hit_rate_redis()
    usage = json.loads(top_q.get())
    pprint(usage)
    top_p.join()
    return res


def run_resource_usage():
    req = 2000
    res = run_once(req, cm="true")
    with open("res.json", "w") as f:
        json.dump(res, f, indent=2)
    del res["raw"]
    pprint(res)


def main():
    reqs = [500, 1000, 1500, 2000, 2500, 3000, 3500, 4000]
    ttls = [100, 1000, 10000]  ## in ms
    uppers_ttl = defaultdict(dict)

    for req in reqs:
        for ttl in ttls:
            upper_ttl = run_once(req, cm="upper", ttl=ttl)
            uppers_ttl[ttl][req] = upper_ttl
    clean2()

    print(uppers_ttl)

    with open(f"{APP}-upper-ttl.json", "w") as f:
        json.dump(uppers_ttl, f, indent=2)


if __name__ == "__main__":
    main()
