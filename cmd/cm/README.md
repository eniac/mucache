## Running and testing the cache manager

To run a scenario with 2 services (with real cache managers and their caches):
```shell
## Better use tmux because we need multiple open connections
# 1: Cache 1
memcached -p 11211 -m 64MB

# 2: Cache 2
memcached -p 11212 -m 64MB

# 3: Cache Manager 1
go build -o cm cmd/cm/main.go ; ./cm --conn http --port 8080 --cache_url localhost:11211 --cm_adds "./experiments/local_cm/cache_manager_addresses.txt"

# 4: Cache Manager 1
go build -o cm cmd/cm/main.go ; ./cm --conn http --port 8081 --cache_url localhost:11212 --cm_adds "./experiments/local_cm/cache_manager_addresses.txt"

# 5: Client
go run cmd/cm/client/main.go --scenario http_two_service
```

To run a simpler http test:
```shell
# On one window
go run cmd/cm/main.go --conn http

# On a second window either do
go run cmd/cm/client/main.go --scenario http
# or do 
curl -X POST localhost:8080/start -d '{"callargs": "req1"}'
curl -X POST localhost:8080/start -d '{"callargs": "req2"}'
curl -X POST localhost:8080/inv -d '{"key": "k1"}'
curl -X POST localhost:8080/end -d '{"callargs": "req1", "caller": "service1", "deps": ["k1", "k2", "k3"], "returnval": "ret1"}'
curl -X POST localhost:8080/end -d '{"callargs": "req2", "caller": "service1", "deps": ["k2", "k3"], "returnval": "ret1"}'
curl -X POST localhost:8080/inv -d '{"key": "k3"}'
```

To run a small test with a memcached server:
```shell
# On one window
memcached -p 11211 -m 64MB

# On another window
go run cmd/cm/client/main.go --scenario test_memcached
```

__NOT SUPPORTED ANYMORE__ To run it with unix socket :
```shell
# On one window
go run cmd/cm/main.go --conn unix_sock

# On a second window
go run cmd/cm/client/main.go --scenario unix_sock
```


## Running a load test

```shell
## You need three windows

## 1st window
go build -o cm cmd/cm/main.go
./cm --conn http

## 2nd window
cd proxy
cargo run --release -- -m cache_manager localhost

## 3rd window
oha -q 20000 -z 10s --latency-correction --disable-keepalive --no-tui http://127.0.0.1:3000
```