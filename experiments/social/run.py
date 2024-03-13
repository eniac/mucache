#!/usr/bin/env python3
import signal
import time
from experiments.helper import *
from pprint import pprint

APP = "social"
set_app(APP)


def start_proxy():
    run_shell("cd proxy && cargo build --release")
    compose_post_ip = get_ip("compose_post")
    home_timeline_ip = get_ip("home_timeline")
    user_timeline_ip = get_ip("user_timeline")
    p = run_in_bg(
        f"cargo run --release social --compose-post {compose_post_ip} --home-timeline {home_timeline_ip} --user-timeline {user_timeline_ip}",
        "proxy")
    time.sleep(5)
    return p


def populate():
    args = ""
    for service in ["social_graph", "compose_post"]:
        ip = get_ip(service)
        args += f" --{service} {ip}"
    args += f" --post_size 100"
    args += f" --number_posts_per_user 10"
    run_shell("python3 experiments/social/populate.py" + args)


def run_once(req: int, cm: str):
    clean2("20")
    deploy(cm=cm)
    populate()
    p = start_proxy()
    top_p, top_q = top_process()
    res = run_shell(compose_oha_proxy(req=req, duration=120))
    res = parse_res(res)
    os.kill(p.pid, signal.SIGINT)
    p.terminate()
    p.wait()
    if cm == 'true' or cm == 'upper':
        res["hit_rate"] = get_hit_rate_redis()
    usage = json.loads(top_q.get())
    pprint(usage)
    top_p.join()
    return res


def run_resource_usage():
    reqs = 1000
    res = run_once(reqs, cm="true")
    print(res['raw'])
    del res['raw']
    pprint(res)


def main():
    reqs = [600, 800, 1000, 1200, 1400, 1600, 1800, 2000]
    baselines = {}
    ours = {}
    uppers = {}

    for req in reqs:
        baseline = run_once(req, cm="false")
        baselines[req] = baseline
        our = run_once(req, cm="true")
        ours[req] = our
        upper = run_once(req, cm="upper")
        uppers[req] = upper
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
