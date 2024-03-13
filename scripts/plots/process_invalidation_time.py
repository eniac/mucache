import argparse

import statistics

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--log_file",
                        help="the log file",
                        type=str,
                        required=True)
    args = parser.parse_args()
    return args

def read_log_file(log_file):
    with open(log_file) as f:
        raw_lines = f.readlines()
    
    lines = [line.rstrip() for line in raw_lines
             if line.startswith("Invalidation")]
    return lines

def main(args):
    lines = read_log_file(args.log_file)

    inv_times = []
    write_diff_times = []
    ## Skip the first line since the experiment hasn't even started
    for line in lines:
        tok1, tok2 = line.split("---")
        inv_time = int(tok1.split(":")[1])
        write_diff_time = int(tok2.split(":")[-1].split(")")[0])
        inv_times.append(inv_time)
        write_diff_times.append(write_diff_time)

    ## ns -> ms
    inv_median = statistics.median(inv_times) / 1_000_000.0
    inv_avg = statistics.mean(inv_times) / 1_000_000.0
    inv_max = max(inv_times) / 1_000_000.0

    write_median = statistics.median(write_diff_times) / 1_000_000.0
    write_avg = statistics.mean(write_diff_times) / 1_000_000.0
    write_max = max(write_diff_times) / 1_000_000.0

    print(f"Invalidation (median: {inv_median}ms) (avg: {inv_avg}ms) (max: {inv_max}ms)")
    print(f"Write diff (median: {write_median}ms) (avg: {write_avg}ms) (max: {write_max}ms)")


if __name__ == '__main__':
    args = parse_arguments()
    main(args)

