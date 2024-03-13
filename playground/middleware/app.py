import os
import uvicorn
import requests
import json
from fastapi import FastAPI, Request
from dapr.ext.fastapi import DaprApp
from dapr.clients import DaprClient

app = FastAPI()
dapr_app = DaprApp(app)

dapr_port = os.getenv("DAPR_HTTP_PORT", 3500)
port = os.getenv("APP_PORT", 3000)

DAPR_STORE_NAME = "statestore"


def http_get_state(key: str):
    try:
        return requests.get(f"http://localhost:{dapr_port}/v1.0/state/{DAPR_STORE_NAME}/{key}").json()
    except json.decoder.JSONDecodeError:
        return None


def http_save_state(key: str, value: str):
    return requests.post(f"http://localhost:{dapr_port}/v1.0/state/{DAPR_STORE_NAME}", json=[{
        "key": key,
        "value": value
    }])


@app.post("/echo")
async def echo(req: Request):
    return await req.body()


@app.post("/state1")
async def invoke(req: Request):
    with DaprClient() as client:
        client.delete_state(DAPR_STORE_NAME, "a")
        client.delete_state(DAPR_STORE_NAME, "A")
        client.save_state(DAPR_STORE_NAME, "a", json.dumps("a"))
        res = client.get_state(DAPR_STORE_NAME, "a").data.decode("utf-8")
        res = json.loads(res)
        assert res == "a", f"Expected a, got {res}"
        res = http_get_state("a")
        assert res == "a", f"Expected a, got {res}"
        res = http_get_state("A")
        assert res is None, f"Expected None, got {res}"
    return "OK"


@app.post("/state2")
async def invoke(req: Request):
    with DaprClient() as client:
        client.delete_state(DAPR_STORE_NAME, "a")
        client.delete_state(DAPR_STORE_NAME, "A")
        http_save_state("a", "a")
        res = client.get_state(DAPR_STORE_NAME, "a").data.decode("utf-8")
        assert res == "", f"Expected '', got {res}"
        res = http_get_state("a")
        assert res is None, f"Expected None, got {res}"
        res = client.get_state(DAPR_STORE_NAME, "A").data.decode("utf-8")
        res = json.loads(res)
        assert res == "A", f"Expected A, got {res}"
        res = http_get_state("A")
        assert res == "A", f"Expected A, got {res}"
    return "OK"


if __name__ == "__main__":
    uvicorn.run(app, host="localhost", port=port)
