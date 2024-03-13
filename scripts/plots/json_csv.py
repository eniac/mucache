import json
import sys

str_data = sys.stdin.read()

data = json.loads(str_data)
column_indexes = sorted([int(key) for key in data.keys()])
rows = ["throughput"] + [f'{i}%' for i in [50, 95]]

for row in rows:
    tokens = []
    for c in column_indexes:
        # print(c)
        tokens.append(str(float(data[str(c)][row])))
    print(",".join(tokens))

# print(column_indexes)

hit_rate_flag = False
for c in column_indexes:
    column_data = data[str(c)]
    if "hit_rate" in column_data:
        hit_rate_flag = True
        hit_rate_keys = list(column_data["hit_rate"].keys())

if hit_rate_flag:
    print("Hit rates")
    print(hit_rate_keys)
    for svc in hit_rate_keys:
        tokens = []
        for c in column_indexes:
            tokens.append(str(float(data[str(c)]["hit_rate"][svc])))
        print(",".join(tokens))
