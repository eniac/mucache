# Mucache Experiments

## Prerequisites

### Cloudlab Cluster
Initialize a cluster of 13 machines on Cloudlab using the following profile:
https://www.cloudlab.us/p/8e054430c6b73652669ec36f24c1ecb716f80c46

Wait until all nodes are ready, download the minifest file as `./manifest.xml`.

### Controller Machine
The controller send files and commands to all worker machines.
Any machine (e.g. Your host) that has Python3 can acts as the controller.

On the controller

```bash
pip3 install fabric
git clone https://github.com/eniac/mucache
```

copy the minifest file to the mucache/

## Setup
On the controller

```bash
cd mucache
export node_username=${username_of_cloudlab}
export private_key=${ssh_key_location}
python3 scripts/host/upload.py
python3 scripts/host/setup.py
```

After it finishes, login into node-0 and you should see the following:

```bash
kubectl get nodes
```

```bash
NAME     STATUS   ROLES           AGE    VERSION
node-0   Ready    control-plane   168m   v1.26.1
node-1   Ready    <none>          168m   v1.26.1
node-2   Ready    <none>          167m   v1.26.1
node-3   Ready    <none>          167m   v1.26.1
node-4   Ready    <none>          167m   v1.26.1
node-5   Ready    <none>          167m   v1.26.1
...
```

## Run
### Build applications
In this artifact, we have four open-source microservice applications, including SocialMedia,
MovieReview, HotelRes, and OnlineBoutique; we have four synthetic applications, including
Proxy, Chain, Fanout, and Fanin.
To run these applications, you can use our [pre-built images on dockerhub](https://hub.docker.com/repository/docker/tauta/mucache/general).

Alternatively, you can build the images on any x86 machines and push to dockerhub by running
```bash
./scripts/host/build_and_push=${docker_io_username} # application image
./scripts/cm/build_and_push=${docker_io_username} # cache manager image
```

If you're using pre-built images, set on the controller node
```bash
export docker_io_username=tauta
```

### Cache Size (Figure 13)
```bash
python3 experiments/cachesize/run.py
```
The results are printed to the console and saved in `hotel-md.json` and `hotel-md-upper.json`. `hotel-md.json` contains a json map with keys being the cache size in MB and values being the p50 and p95 latency. `hotel-md-upper` is the TTL-inf baseline.

### Microbenchmark (Figure 14)
To run the three microservices, Chain, Fanout and Fanin.
```bash
python3 experiments/chain/run.py # chain
python3 experiments/star/run.py  # fanout
python3 experiments/fanin/run.py # fanin
```

The results will be saved to `{APP}-baseline.json` for baseline and `{APP}.json` for mucache.

### Hitrate (Figure 17)
```bash
python3 experiments/twoservices/hitrate.py
```
The results are saved at `hitrate_hdr_baseline` for the baseline, `hitrate_hdr_{hit_rate}` for mucache and `hitrate_hdr_{hit_rate}_upper` for TTL-inf.

### Real-world applications (Figure 10)
To run the four real-world applications, run
```bash
python3 experiments/movie/run.py
python3 experiments/hotel/run.py
python3 experiments/boutique/run.py
python3 experiments/social/run.py
```
The scripts will print the output and write the results to `{APP}-upper.json`, `{APP}-baseline.json` and `{APP}.json`, which correspond to TTL-inf, baseline and mucache respectively.

### Different TTL baselines (Figure 11)

```bash
python3 experiments/hotel/ttl.py
```
The script will print the output and write and result to `hotel-upper-ttl.json`.

### Sharding (Figure 12)
To build the applications and cache manager for the sharding example, run

```bash
./scripts/host/shard_build_and_push=${docker_io_username} # application image
./scripts/cm/shard_build_and_push=${docker_io_username} # cache manager image
```