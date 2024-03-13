# golang

## Servers

```bash
## Backend
go build -o main main.go
dapr run --app-id goapp --app-port 3000 --dapr-http-port 3500 ./main --config ./components/config.yaml --components-path ./components/

## Frontend
cd frontend
go build -o frontend main.go
dapr run --app-id frontend --app-port 3001 --dapr-http-port 3501 ./frontend --config ./../components/config.yaml --components-path ./../components/
```

## Client
    
```bash
curl http://localhost:3500/v1.0/invoke/goapp/method/echo
~/wrk2/wrk -t2 -c16 -d20s -R10000 --latency http://localhost:3500/v1.0/invoke/goapp/method/echo


curl -X POST -s http://localhost:3500/v1.0/invoke/goapp/method/getset -H Content-Type:application/json --data @reservation.json
./wrk2/wrk -t1 -c1 -d20 -R1 --latency http://localhost:3500/v1.0/invoke/goapp/method/getset -s getset-workload.lua


curl -X POST -s http://localhost:3501/v1.0/invoke/frontend/method/getset -H Content-Type:application/json --data @reservation.json
./wrk2/wrk -t1 -c1 -d20 -R1 --latency http://localhost:3501/v1.0/invoke/frontend/method/getset -s getset-workload.lua
```

## Measurements

On m510 cloudlab

Backend (no work):
- 1 Worker
    - Max throughput about 10k RPS (wrk -t2 -c16) 

Backend (GET-SET):
- 1 Worker 
    - Max throughput about 4k RPS (wrk -t2 -c16) (goserver 100% CPU, daprd 550% CPU)
    - Latency 5 - 5.6ms (50th - 95th) (wrk -t1 -c1 -R10)

2 Services:
- 1 worker each
    - Max throughput about 2k RPS (wrk -t2 -c16) (backend 90% CPU, daprd-backend 550% CPU, daprd-frontend 180%, frontend 70% -- total almost 850%)
    - Latency 7.7 - 9.3ms (50th - 95th) (wrk -t1 -c1 -R10)
