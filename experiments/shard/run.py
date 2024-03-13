#!/usr/bin/env python3

import signal
import time
from experiments.helper import *

APP = "social"
set_app(APP)


# def generate_cm_adds(shard_count: int):
#     with open(f"{project_path()}/experiments/k8s_cm/social-shard.txt", "w") as f:
#         for idx, service in enumerate(APPS_NO_UNDERSCORE[APP]):
#             for shard in range(1, shard_count + 1):
#                 f.write(f"{service}{shard} http://cm{idx + 1}{shard}\n")


def start_proxy():
    run_shell("cd proxy && cargo build --release")
    compose_post_ip = get_shard_ips("compose_post")
    home_timeline_ip = get_shard_ips("home_timeline")
    user_timeline_ip = get_shard_ips("user_timeline")
    p = run_in_bg(
        f"cargo run --release social --compose-post {compose_post_ip} --home-timeline {home_timeline_ip} --user-timeline {user_timeline_ip}",
        "proxy")
    time.sleep(5)
    return p


def populate():
    args = ""
    for service in ["social_graph", "compose_post"]:
        ip = get_ip(service + "1")
        args += f" --{service} {ip}"
    args += f" --post_size 100"
    args += f" --number_posts_per_user 10"
    cmd = f"python3 experiments/social/populate.py" + args
    print(cmd)
    run_shell(cmd)


def run_once(shard_count: int, req: int, max_shard_count: int, cm: bool):
    # generate_cm_adds(shard_count)
    shard_clean2(max_shard_count)
    shard_deploy(shard_count, cm=cm)
    populate()
    p = start_proxy()
    res = run_shell(compose_oha_n_proxy(req))
    res = parse_res(res)
    os.kill(p.pid, signal.SIGINT)
    p.terminate()
    p.wait()
    if cm:
        res["hit_rate"] = get_hit_rate_redis()
    return res


def main():
    shard_counts = [1, 2, 4]
    max_shard_count = max(shard_counts)
    max_shard_count = 4
    res = {}
    baseline = {}

    for shard_count in shard_counts:
        # shard_res = {}
        # for req in reqs[shard_count]:
        #     shard_res[req] = run_once(shard_count, req, max_shard_count, True)
        # res[shard_count] = shard_res
        res[shard_count] = run_once(shard_count, 50000, max_shard_count, True)
        # shard_res = {}
        # for req in baseline_reqs[shard_count]:
        #     shard_res[req] = run_once(shard_count, req, max_shard_count, False)
        # baseline[shard_count] = shard_res
        baseline[shard_count] = run_once(shard_count, 50000, max_shard_count, False)

    print(res)
    print(baseline)
    with open(f"{APP}-shard.json", "w") as f:
        json.dump(res, f, indent=2)
    with open(f"{APP}-shard-baseline.json", "w") as f:
        json.dump(baseline, f, indent=2)


if __name__ == "__main__":
    main()
