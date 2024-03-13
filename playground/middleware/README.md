# Middleware Playground

## Setup
Update dapr config to config.yaml

Add uppercase.yaml to components

## Server

```bash
# Server
# pwd: playground/middleware
dapr run --app-id pythonapp --app-port 3000 --dapr-http-port 3500 python3 app.py
```


## Client

### HTTP Interceptor
```bash
# Client
curl -XPOST -d hello http://localhost:3500/v1.0/invoke/pythonapp/method/echo
# should return "HELLO"
```

### State Interceptor
```bash
curl -XPOST -d hello http://localhost:3500/v1.0/invoke/pythonapp/method/state1
curl -XPOST -d hello http://localhost:3500/v1.0/invoke/pythonapp/method/state2
# both should return "OK"
```

#### Reads
Reads are As-Is

| Client\Database       | a: a | A: A |
|-----------------------|------|------|
| Dapr Client (get a)   | a    | None |
| Dapr Client (get A)   | None | A    |
| HTTP requests (get a) | a    | None |
| HTTP requests (get A) | None | A    |

#### Writes
HTTP writes goes through middlewares

|                            | Database |
|----------------------------|----------|
| Dapr Client (write a: a)   | a: a     |
| HTTP requests (write a: a) | A: A     |


## Custom Middleware

### Compile dapr
```bash
# pwd: dapr
make build
cp ./dist/linux_amd64/release/daprd ~/.dapr/bin
```