import argparse

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
             if line.startswith("Processed")]
    return lines

def main(args):
    lines = read_log_file(args.log_file)

    total_items = 0
    total_time_ms = 0
    ## Skip the first line since the experiment hasn't even started
    for line in lines[1:]:
        items = int(line.split()[1])
        time_ms = int(line.split()[-1])
        total_items += items
        total_time_ms += time_ms

    total_time_s = total_time_ms / 1_000_000.0    
    items_per_second = total_items / total_time_s
    print("Cache manager events per second:", items_per_second)


if __name__ == '__main__':
    args = parse_arguments()
    main(args)

