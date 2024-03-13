
from app import User

import json
import logging
from typing import Union

from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

from dapr.clients import DaprClient

class Reservation(BaseModel):
    user_name: str
    user_id: str
    hotel_id: str


#code
DAPR_STORE_NAME = "statestore"
logging.basicConfig(level = logging.INFO)

@app.post("/book")
async def book(reservation: Reservation):
    ## TODO: Can we not restart the client on every request?
    with DaprClient() as client:
        #Using Dapr SDK to save and get state
        user = User(name=reservation.user_name, id=reservation.user_id)

        ## NOTE: It is important to use the async invocation API,
        ##       because the standard one leaks open files.
        resp = await client.invoke_method_async(app_id="backend",
                                                method_name=f"/book_hotel/{reservation.hotel_id}",
                                                data=user.json(),
                                                http_verb="POST",
                                                content_type="application/json")
        
        ## TODO: This load here is unnecessary (since we will
        ##       need to reserialize it here). Can we just return the string?
        ##       
        ##       FastAPI has type for request (and does auto encode/decode)
        ##       but not for response!
        json_resp = json.loads(resp.data)
        # logging.info('Data returned: ' + str(resp.data))

    return json_resp



# app.run(host="0.0.0.0",port=appPort)
