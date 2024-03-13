#!/usr/bin/env python3
import signal
import time
from experiments.helper import *
from pprint import pprint

APP = "boutique"
set_app(APP)

## Increase this to increase the latency of home requests (10% of requests)
CATALOG_SIZE = 10

## Tweaking this will vary how of1ten the cart of a user is accessed
USERS = 10000

## Tweaking this will vary how many objects there are and their size
PRODUCTS = 10000
PRODUCT_SIZE = 10000

## Tweak the cache size for the hit ratio (in MB)
CACHE_SIZE = "80"
# CACHE_SIZE = "0"

## Tweaking the TTL affects the TTL baseline (in ms)
TTL = "10"


def start_proxy():
    global CATALOG_SIZE
    run_shell("cd proxy && cargo build --release")
    frontend_ip = get_ip("frontend")
    cart_ip = get_ip("cart")
    currency_ip = get_ip("currency")
    p = run_in_bg(
        f"cargo run --release boutique --frontend {frontend_ip} --cart {cart_ip} --currency {currency_ip} --catalog-size {CATALOG_SIZE}",
        "proxy")
    time.sleep(5)
    return p


def populate():
    global USERS, PRODUCTS, PRODUCT_SIZE
    args = ""
    args += f" --users {USERS}"
    args += f" --products {PRODUCTS}"
    args += f" --product_size {PRODUCT_SIZE}"
    for service in ["frontend", "product_catalog", "currency"]:
        ip = get_ip(service)
        args += f" --{service} {ip}"
    run_shell("python3 experiments/boutique/populate.py" + args)


def run_once(req: int, cm: str, ttl=None):
    global CACHE_SIZE
    clean2(mem=CACHE_SIZE)
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
    reqs = 3500
    res = run_once(reqs, cm="true")
    print(res['raw'])
    del res['raw']
    pprint(res)


def main():
    reqs = [2000, 3000, 4000, 4500, 5000, 5500, 6000]
    reqs = [5000, 5500, 6000]
    ttl = TTL  ## in ms
    baselines = {}
    uppers = {}
    ours = {}

    ## Note: Save every iteration so that we get incremental results
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

if __name__ == "__main__":
    main()
