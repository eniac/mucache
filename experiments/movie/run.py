#!/usr/bin/env python3
import signal
import time
from experiments.helper import *
from pprint import pprint

APP = "movie"
set_app(APP)


def start_proxy():
    run_shell("cd proxy && cargo build --release")
    frontend_ip = get_ip("frontend")
    page_ip = get_ip("page")
    p = run_in_bg(
        f"cargo run --release movie --frontend {frontend_ip} --page {page_ip}",
        "proxy")
    time.sleep(5)
    return p


def populate():
    args = ""
    for service in APPS[APP]:
        ip = get_ip(service)
        args += f" --{service} {ip}"
    run_shell("python3 experiments/movie/populate.py" + args)


def run_once(req: int, cm: str, ttl=None):
    clean2(mem="70")
    deploy(cm=cm, ttl=ttl)
    populate()
    p = start_proxy()
    top_p, top_q = top_process()
    res = run_shell(compose_oha_proxy(req=req))
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
    req = 2500
    res = run_once(req, cm="true")
    with open("res.json", "w") as f:
        json.dump(res, f, indent=2)
    del res["raw"]
    pprint(res)


def main():
    reqs = [1000, 1500, 2000, 2500, 3000, 3500, 4000]
    ttl = "6"  ## in ms
    baselines = {}
    uppers = {}
    ours = {}

    for req in reqs:
        baseline = run_once(req, cm="false")
        baselines[req] = baseline
        with open(f"{APP}-baseline.json", "w") as f:
            json.dump(baselines, f, indent=2)

        upper = run_once(req, cm="upper")
        uppers[req] = upper
        with open(f"{APP}-upper.json", "w") as f:
            json.dump(uppers, f, indent=2)

        our = run_once(req, cm="true")
        ours[req] = our
        with open(f"{APP}.json", "w") as f:
            json.dump(ours, f, indent=2)
    clean2()

    print(baselines)
    print(ours)
    print(uppers)
    with open(f"{APP}-baseline.json", "w") as f:
        json.dump(baselines, f, indent=2)
    with open(f"{APP}.json", "w") as f:
        json.dump(ours, f, indent=2)
    with open(f"{APP}-upper.json", "w") as f:
        json.dump(uppers, f, indent=2)


if __name__ == "__main__":
    main()
