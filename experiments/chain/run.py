#!/usr/bin/env python3
import signal
import time
from experiments.helper import *

APP = "chain"
set_app(APP)


# def start_proxy():
#     run_shell("cd proxy && cargo build --release")
#     compose_post_ip = get_ip("compose_post")
#     home_timeline_ip = get_ip("home_timeline")
#     user_timeline_ip = get_ip("user_timeline")
#     p = run_in_bg(
#         f"cargo run --release social --compose-post {compose_post_ip} --home-timeline {home_timeline_ip} --user-timeline {user_timeline_ip}",
#         "proxy")
#     time.sleep(5)
#     return p




def run_once(req: int, cm: bool, endpoint_suffix: str, data):
    clean()
    deploy(cm=cm)
    frontend_ip = get_ip("service1")
    endpoint= f'http://{frontend_ip}/{endpoint_suffix}'
    # p = start_proxy()
    res = run_shell(compose_oha_proxy_post(data, endpoint=endpoint, req=req))
    # print("oha results")
    # print(res)
    res = parse_res(res)
    # os.kill(p.pid, signal.SIGINT)
    # p.terminate()
    # p.wait()
    if cm:
        res["hit_rate"] = get_all_hit_rate()
    return res


def main():
    reqs = [1000, 2000, 3000, 4000, 5000, 6000]
    baselines = {}
    ours = {}


    ## TODO: We might want some more complex input workload,
    ##       with a distribution and writes
    baseline_endpoint = f'/ro_read'
    baseline_data = {
        "k": 1
    }

    ours_endpoint = f'/ro_hitormiss'
    ours_data = {
        "k": 1,
        "hit_rate": 0.5
    }

    for req in reqs:
        baseline = run_once(req, cm=False, endpoint_suffix=baseline_endpoint, data=baseline_data)
        print("Baseline:", baseline)
        baselines[req] = baseline
        our = run_once(req, cm=True, endpoint_suffix=ours_endpoint, data=ours_data)
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
