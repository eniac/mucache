import json
from pprint import pprint
from collections import defaultdict


def parse_ttl():
    res = {}
    with open("../hotel-upper-ttl.json") as f:
        res = json.load(f)
    for ttl, vs in res.items():
        rpss = []
        p50s = []
        p95s = []
        throughputs = []
        hit_rates = defaultdict(list)
        for rps, v in vs.items():
            rpss.append(rps)
            p50s.append(str(v["50%"]))
            p95s.append(str(v["95%"]))
            throughputs.append(str(v["throughput"]))
            for service, hit_rate in v["hit_rate"].items():
                hit_rates[service].append(str(hit_rate))
        print(ttl)
        print("\t".join(rpss))
        print("\t".join(throughputs))
        print("\t".join(p50s))
        print("\t".join(p95s))
        for service, hit_rate in hit_rates.items():
            print(service)
            print("\t".join(hit_rate))


def parse():
    ours = {}
    with open("../social-md.json") as f:
        ours = json.load(f)
    rps = []
    p50s = []
    p95s = []
    for k, v in ours.items():
        rps.append(k)
        p50s.append(str(v["50%"]))
        p95s.append(str(v["95%"]))
    print("\t".join(rps))
    print("\t".join(p50s))
    print("\t".join(p95s))


def parse_md():
    ours = {}
    with open("../social-md.json") as f:
        ours = json.load(f)
    mems = []
    p50s = []
    p95s = []
    hit_rates = defaultdict(dict)
    for k, v in ours.items():
        mems.append(k)
        p50s.append(str(v["50%"]))
        p95s.append(str(v["95%"]))
        for service, hit_rate in v["hit_rate"].items():
            hit_rates[service][k] = str(hit_rate)
    print("\t".join(mems))
    print("\t".join(p50s))
    print("\t".join(p95s))
    for service, hit_rate in hit_rates.items():
        print(service, end="\t")
        print("\t".join(hit_rate.values()))


def parse_hitrate():
    ours = {}
    with open("../hitrate-upper.json") as f:
        ours = json.load(f)
    hitrates = []
    p50s = []
    p95s = []
    for k, v in ours.items():
        hitrates.append(k)
        p50s.append(str(v["50%"]))
        p95s.append(str(v["95%"]))
    print("\t".join(hitrates))
    print("\t".join(p50s))
    print("\t".join(p95s))


def main():
    # parse()
    # parse_ttl()
    # parse_md()
    parse_hitrate()


if __name__ == '__main__':
    main()
