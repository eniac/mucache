To install fastAPI we need both fastAPI and uvicorn:
```sh
python3 -m pip install fastapi "uvicorn[standard]"
```

Run the two services by running the following two commands (in different windows)

```sh
## Backend
dapr run --log-level error --app-id backend --app-port 5001 --dapr-http-port 3501 bash run_server.sh 1 critical

## Frontend
dapr run --log-level error --app-id frontend --app-port 5002 --dapr-http-port 3502 bash run_frontend.sh 1 critical
```

And then send requests with the following:

```sh
# Frontend
curl -X POST -s http://localhost:5002/book -H Content-Type:application/json --data @reservation.json

# Backend
curl -X POST -s http://localhost:5001/book_hotel/10 -H Content-Type:application/json --data @user.json
```

Send load to the application using (after having installed wrk2 in the current directory):
```sh
./wrk2/wrk -t1 -c1 -d20 -R1 --latency http://localhost:5002/book -s workload.lua
```

## Performance characteristics

Remember to set logging to minimal and modify the config to have minimal tracing and disabled metrics:
```yaml
spec:
  tracing:
    samplingRate: "0"
    zipkin:
      endpointAddress: http://localhost:9411/api/v2/spans
  metrics:
    enabled: false
```


On m510 cloudlab:

Backend without any DB access and without initializing the Dapr Client:
- 1 Worker
    - Max throughput about 3.3k RPS (wrk -t4 -c100) (~100% CPU on the uvicorn server)
    - Latency 2.1ms - 2.9ms (50th - 95th) (wrk -t1 -c1 -R10)
    - Numbers without Dapr don't change a lot. It seems that Dapr adds some overhead in cross-service calls but not only in the backend

Backend without any DB access (just initializing the Dapr Client):
- 1 Worker
    - Max throughput about 600 RPS (wrk -t2 -c50) (~100% CPU on the uvicorn server)
    - Latency 3.5ms - 4.4ms (50th - 95th) (wrk -t1 -c1 -R10)

Backend:
- 1 Worker
    - Max throughput about 200 RPS (wrk -t2 -c50) (about 80% CPU on the uvicorn server)
    - Latency 8.6ms - 9.5ms (50th - 95th) (wrk -t1 -c1 -R10)
- 2 workers
    - Max throughput about 380 RPS (wrk -t2 -c70) (about 80% CPU on the uvicorn server)

2 Services (Backend without any DB access and without initializing the Dapr Client):
- 1 worker each
    - Max throughput about 450 RPS (wrk -t2 -c50)
    - Latency 9.8ms - 10.8ms (50th - 95th) (wrk -t1 -c1 -R10)
- 2 workers frontend
    - Max throughput about 500 RPS (wrk -t2 -c70)
- 4 workers frontend
    - Max throughput about 700 RPS (wrk -t2 -c70)


2 Services:
- 1 worker each
    - Max throughput about 150 RPS (wrk -t2 -c50)
    - Latency 16ms - 19ms (50th - 95th) (wrk -t1 -c1 -R10)
- 2 workers each
    - Max throughput about 300 RPS (wrk -t2 -c70)


## TODO Items

__TODO:__ Figure out how to clean the state store.

__TODO:__ Check out distributed locks (because they do not support transactional updates).