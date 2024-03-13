import json
import time
from collections import defaultdict
from multiprocessing import Process, Queue

from scripts.host.common import *

NWORKERS = len(SERVERS) - 1
APP = ""
DOCKER_NAME = os.getenv("docker_io_username")
if DOCKER_NAME is None:
    exit("Please set docker_io_username")

APPS = {
    "twoservices": ["caller", "callee"],
    "chain": ["service1", "service2", "service3", "service4", "backend"],
    "star": ["frontend", "backend1", "backend2", "backend3", "backend4"],
    "fanin": ["frontend1", "frontend2", "frontend3", "frontend4", "backend"],
    "social": ["post_storage", "home_timeline", "user_timeline", "social_graph", "compose_post"],
    "hotel": ["frontend", "profile", "rate", "reservation", "search", "user"],
    "boutique": ["cart", "checkout", "currency", "email", "frontend", "payment", "product_catalog", "recommendations",
                 "shipping"],
    "movie": ["cast_info", "compose_review", "frontend", "movie_id", "movie_info", "movie_reviews", "page", "plot",
              "review_storage", "unique_id", "user", "user_reviews"],
}
APPS_NO_UNDERSCORE = {}

for k, v in APPS.items():
    APPS_NO_UNDERSCORE[k] = ["".join(s.split("_")) for s in v]

VALID_CM_CHOICES = ["false", "true", "upper"]


def set_app(app: str):
    global APP
    APP = app


def get_ip(service: str):
    return run_shell(f'bash -c "source {PROJECT_PATH}/scripts/utility.sh && kip {service}"').strip()


def get_shard_ips(service: str):
    return run_shell(f'bash -c "source {PROJECT_PATH}/scripts/utility.sh && kips {service}"').strip()


def get_hit_rate(service: str) -> float:
    assert False, "Deprecated"
    hit = run_shell(f'bash -c "source {PROJECT_PATH}/scripts/utility.sh && klog {service} | grep hit | wc -l"').strip()
    miss = run_shell(
        f'bash -c "source {PROJECT_PATH}/scripts/utility.sh && klog {service} | grep miss | wc -l"').strip()
    hit = float(hit)
    miss = float(miss)
    if hit + miss == 0:
        return 0
    return hit / (hit + miss)


def get_all_hit_rate() -> dict:
    assert False, "Deprecated"
    res = {}
    for service in APPS[APP]:
        res[service] = "{:.2f}".format(get_hit_rate(service))
    return res


def get_hit_rate_redis() -> dict:
    res = {}
    for idx in range(len(APPS[APP])):
        service = APPS[APP][idx]
        ip = get_ip(f"cache{idx + 1}-redis-master")
        hit = run_shell(f'echo "info stats" | redis-cli -h {ip} | grep keyspace_hits').split(":")[1].strip()
        miss = run_shell(f'echo "info stats" | redis-cli -h {ip} | grep keyspace_misses').split(":")[1].strip()
        hit = float(hit)
        miss = float(miss)
        hitrate = 0
        if hit + miss != 0:
            hitrate = hit / (hit + miss)
        if hitrate > 0:
            res[service] = "{:.2f}".format(hitrate)
    return res


def get_hit_miss_redis() -> dict:
    res = {}
    for idx in range(len(APPS[APP])):
        service = APPS[APP][idx]
        ip = get_ip(f"cache{idx + 1}-redis-master")
        hit = run_shell(f'echo "info stats" | redis-cli -h {ip} | grep keyspace_hits').split(":")[1].strip()
        miss = run_shell(f'echo "info stats" | redis-cli -h {ip} | grep keyspace_misses').split(":")[1].strip()
        hit = float(hit)
        miss = float(miss)
        res[service] = (hit, miss)
    return res


def compute_hit_rate_redis(old: dict, new: dict) -> dict:
    res = {}
    for service in old:
        old_hit, old_miss = old[service]
        new_hit, new_miss = new[service]
        hitrate = 0
        if new_hit - old_hit + new_miss - old_miss != 0:
            hitrate = (new_hit - old_hit) / (new_hit - old_hit + new_miss - old_miss)
        if hitrate > 0:
            res[service] = "{:.2f}".format(hitrate)
    return res


def clean(mem=None):
    assert False, "Deprecated"
    assert APP in APPS, f"Invalid app: {APP}"
    print("Cleaning up...")
    run_shell(f"{PROJECT_PATH}/scripts/clean.sh {APP}")
    if mem is not None:
        run_shell(f"{PROJECT_PATH}/scripts/setup/restart_memcached.sh {NWORKERS} {mem}")
    else:
        run_shell(f"{PROJECT_PATH}/scripts/setup/restart_memcached.sh {NWORKERS}")
    run_shell(f"{PROJECT_PATH}/scripts/setup/restart_redis.sh {NWORKERS}")


def clean2(mem=0):
    assert APP in APPS, f"Invalid app: {APP}"
    print("Cleaning up...")
    run_shell(f"{PROJECT_PATH}/scripts/clean.sh {APP}")
    run_shell(f"{PROJECT_PATH}/scripts/setup/restart_cache.sh {len(APPS[APP])} {mem}")
    run_shell(f"{PROJECT_PATH}/scripts/setup/restart_redis.sh {len(APPS[APP])}")
    for i in range(len(APPS[APP])):
        run_shell(f"kubectl rollout status deployment/cache{i + 1}-redis-master")


def shard_clean2(shard_count: int, mem=0):
    assert APP in APPS, f"Invalid app: {APP}"
    print("Cleaning up...")
    run_shell(f"{PROJECT_PATH}/scripts/shard_clean.sh {APP} {shard_count}")
    run_shell(f"{PROJECT_PATH}/scripts/setup/restart_redis.sh {len(APPS[APP])}")
    run_shell(f"{PROJECT_PATH}/scripts/setup/restart_cache.sh {len(APPS[APP])} {mem}")
    for i in range(len(APPS[APP])):
        run_shell(f"kubectl rollout status deployment/cache{i + 1}-redis-master")


def deploy(cm=str, ttl=None):
    print("Deploying...")
    assert (cm in VALID_CM_CHOICES)
    if not ttl is None:
        print(run_shell(f"{PROJECT_PATH}/scripts/deploy.sh {DOCKER_NAME} {APP} {cm} {ttl}"))
    else:
        print(run_shell(f"{PROJECT_PATH}/scripts/deploy.sh {DOCKER_NAME} {APP} {cm}"))
    for service in APPS_NO_UNDERSCORE[APP]:
        run_shell(f"kubectl rollout status deployment/{service}")
    time.sleep(10)
    print("Deployed")


def shard_deploy(shard_count: int, cm=True):
    print("Deploying...")
    if cm:
        print(run_shell(f"{PROJECT_PATH}/scripts/shard_deploy.sh {DOCKER_NAME} {APP} true {shard_count}"))
    else:
        print(run_shell(f"{PROJECT_PATH}/scripts/shard_deploy.sh {DOCKER_NAME} {APP} false {shard_count}"))
    for service in APPS_NO_UNDERSCORE[APP]:
        for i in range(shard_count):
            run_shell(f"kubectl rollout status deployment/{service}{i + 1}")
    print("Deployed")


def compose_oha_cmd(service: str, method: str, data, req: int = 1000, duration: int = 30):
    ip = get_ip(service)
    return f"oha -q {req} -z {duration}s --latency-correction --no-tui http://{ip}/{method} -d '{json.dumps(data)}'"


def compose_wrk_cmd(service: str, method: str, script, req: int = 1000, duration: int = 30):
    ip = get_ip(service)
    return f"$HOME/wrk2/wrk -t4 -c32 -d{duration}s -R{req} -s {script} --latency http://{ip}/{method}"


def compose_oha_proxy(req: int = 1000, duration: int = 30):
    return f"oha -q {req} -z {duration}s --latency-correction --no-tui http://127.0.0.1:3000"


def compose_oha_n_proxy(nreq: int = 10000):
    return f"oha -n {nreq} --latency-correction --no-tui http://127.0.0.1:3000"


def compose_oha_proxy_post(data, endpoint="http://127.0.0.1:3000", req: int = 1000, duration: int = 30):
    oha_invocation_string = f"oha -q {req} -z {duration}s -m POST -H 'Content-Type: application/json' -d '{json.dumps(data)}' --latency-correction --disable-keepalive --no-tui {endpoint}"
    print("Invoking oha")
    print(oha_invocation_string)
    return oha_invocation_string


def compose_curl_cmd(service: str, method: str, data):
    ip = get_ip(service)
    return f"curl -X POST {ip}/{method} -d '{json.dumps(data)}'"


def compose_curl_proxy():
    return f"curl 127.0.0.1:3000"


def parse_res(res, percents=["50%", "95%"]):
    lats = {}
    lats["raw"] = res
    for line in res.splitlines():
        if "Requests/sec" in line:
            throughput = float(line.split()[1])
            lats["throughput"] = throughput
            continue
        if "% in " not in line:
            continue
        ls = line.split()
        in_idx = ls.index("in")
        percent, lat = ls[in_idx - 1], float(ls[in_idx + 1]) * 1000
        if percent in percents:
            lats[percent] = lat
    return lats


def top(q):
    usage = defaultdict(lambda: {
        "cpu": [],
        "mem": [],
    })
    for _ in range(10):
        res = run_shell("kubectl top pods --containers")
        for line in res.splitlines():
            if "NAME" in line or "redis" in line:
                continue
            ls = line.split()
            pod = ls[0].split("-")[0]
            name = ls[1]
            cpu = ls[2]
            mem = ls[3]
            cpu = int(cpu[:-1])
            mem = int(mem[:-2])
            assert pod in APPS_NO_UNDERSCORE[APP] or "cm" in pod
            usage[f"{pod}-{name}"]["cpu"].append(cpu)
            usage[f"{pod}-{name}"]["mem"].append(mem)
        time.sleep(10)
    q.put(json.dumps(usage))


def top_process():
    q = Queue()
    p = Process(target=top, args=(q,))
    p.start()
    return p, q
