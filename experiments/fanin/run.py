#!/usr/bin/env python3
import signal
import time
from experiments.helper import *

APP = "fanin"
set_app(APP)


def start_proxy(endpoint):
    run_shell("cd proxy && cargo build --release")
    frontend1_ip = get_ip("frontend1")
    frontend2_ip = get_ip("frontend2")
    frontend3_ip = get_ip("frontend3")
    frontend4_ip = get_ip("frontend4")
    hitrate = 0.5
    endpoint = endpoint
    p = run_in_bg(
        f"cargo run --release fanin --frontend1 {frontend1_ip} --frontend2 {frontend2_ip} --frontend3 {frontend3_ip} --frontend4 {frontend4_ip} --hitrate {hitrate} --endpoint {endpoint}",
        "proxy")
    time.sleep(5)
    return p




def run_once(req: int, cm: bool, endpoint_suffix: str):
    clean()
    deploy(cm=cm)
    p = start_proxy(endpoint=endpoint_suffix)
    res = run_shell(compose_oha_proxy(req=req))
    # print("oha results")
    # print(res)
    res = parse_res(res)
    os.kill(p.pid, signal.SIGINT)
    p.terminate()
    p.wait()
    if cm:
        res["hit_rate"] = get_all_hit_rate()
    return res


def main():
    reqs = [2000, 4000, 6000, 8000, 10000, 12000, 14000, 16000]
    baselines = {}
    ours = {}

    baseline_endpoint = f'ro_read'
    ours_endpoint = f'ro_hitormiss'

    for req in reqs:
        baseline = run_once(req, cm=False, endpoint_suffix=baseline_endpoint)
        print("Baseline:", baseline)
        baselines[req] = baseline
        our = run_once(req, cm=True, endpoint_suffix=ours_endpoint)
        print("Ours:", our)
        ours[req] = our
    clean2()
    print(baselines)
    print(ours)
    with open(f"{APP}-baseline.json", "w") as f:
        json.dump(baselines, f, indent=2)
    with open(f"{APP}.json", "w") as f:
        json.dump(ours, f, indent=2)


if __name__ == "__main__":
    main()
