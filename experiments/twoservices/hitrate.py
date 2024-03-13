#!/usr/bin/env python3

from experiments.helper import *

APP = "twoservices"
set_app(APP)


def populate():
    data = {
        "k": 1,
        "v": 1,
    }
    populate_res = run_shell(compose_curl_cmd("caller", "write", data))
    assert "OK" in populate_res


def run_baseline():
    clean2()
    deploy(cm="false")
    populate()
    data = {
        "k": 1,
    }
    res = run_shell(compose_oha_cmd("caller", "read", data, req=1000, duration=60))
    res = parse_res(res)
    del res["raw"]
    from pprint import pprint
    pprint(res)
    return res
    # with open(f"hitrate_baseline.json", "w") as f:
    #     json.dump(res, f, indent=2)


def run_once(hit_rate: float, cm: str = "true"):
    clean2()
    deploy(cm=cm)
    populate()
    data = {
        "k": 1,
        "hit_rate": hit_rate,
    }
    res = run_shell(compose_oha_cmd("caller", "ro_hitormiss", data, req=1000, duration=60))
    res = parse_res(res)
    return res
    # with open(f"hitrate_{hit_rate}_{cm}.json", "w") as f:
    #     json.dump(res, f, indent=2)


def run():
    baseline = {}
    ours = {}
    uppers = {}
    baseline["baseline"] = run_baseline()
    for hit_rate in [0.0, 0.2, 0.4, 0.6, 0.8, 1.0]:
        our = run_once(hit_rate)
        upper = run_once(hit_rate, cm="upper")
        ours[hit_rate] = our
        uppers[hit_rate] = upper

    print(baseline)
    print(ours)
    print(uppers)
    with open(f"hitrate-baseline.json", "w") as f:
        json.dump(baseline, f, indent=2)
    with open(f"hitrate.json", "w") as f:
        json.dump(ours, f, indent=2)
    with open(f"hitrate-upper.json", "w") as f:
        json.dump(uppers, f, indent=2)


def run_hdr_once(hit_rate: float, cm: str = "true"):
    clean2()
    deploy(cm=cm)
    populate()
    run_shell(
        f'HITRATE="{hit_rate}" /usr/bin/envsubst < {PROJECT_PATH}/experiments/twoservices/hitrate_template.lua > {PROJECT_PATH}/experiments/twoservices/hitrate.lua')
    res = run_shell(
        compose_wrk_cmd("caller", "ro_hitormiss", f"{PROJECT_PATH}/experiments/twoservices/hitrate.lua", req=1000,
                        duration=120))
    return res


def run_baseline_hdr():
    clean2()
    deploy(cm="false")
    populate()
    run_shell(
        f'/usr/bin/envsubst < {PROJECT_PATH}/experiments/twoservices/hitrate_baseline_template.lua > {PROJECT_PATH}/experiments/twoservices/hitrate_baseline.lua')
    res = run_shell(
        compose_wrk_cmd("caller", "read", f"{PROJECT_PATH}/experiments/twoservices/hitrate_baseline.lua", req=1000,
                        duration=120))
    return res


def run_hdr():
    hit_rates = [0.0, 0.6]
    baseline = run_baseline_hdr()
    with open(f"hitrate_hdr_baseline", "w") as f:
        f.write(baseline)
    for hit_rate in hit_rates:
        our = run_hdr_once(hit_rate)
        with open(f"hitrate_hdr_{hit_rate}", "w") as f:
            f.write(our)
        upper = run_hdr_once(hit_rate, cm="upper")
        with open(f"hitrate_hdr_{hit_rate}_upper", "w") as f:
            f.write(upper)
    clean2()


def main():
    run_hdr()
    # run_baseline()
    # run()


if __name__ == "__main__":
    main()
