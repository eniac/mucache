import logging
from typing import Union

from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

from dapr.clients import DaprClient

class User(BaseModel):
    name: str
    id: str
    is_offer: Union[bool, None] = None


#code
DAPR_STORE_NAME = "statestore"
logging.basicConfig(level = logging.INFO)


@app.get("/")
async def read_root():
    return {"Hello": "World"}

@app.post("/book_hotel/{hotel_id}")
async def read_item(hotel_id: int, user: User):
    ## TODO: Can we not restart the client on every request?
    # with DaprClient() as client:
        # hotel_key = f"hotel_{hotel_id}"
        # ## TODO: Later we should have our own binary encoding of arbitrary 
        # ##       data structures (for now we can just use strings)
        # ## TODO: Check fastapi as it has encoders

        # ## TODO: Provide proper consistency guarantees
        # ##       (1) Optimistic with etags and retries
        # ##       (2) Pessimistic with distributed lock

        # init_value = ""
        # response = client.get_state(DAPR_STORE_NAME, hotel_key)
        # hotel_residents = response.data.decode('utf-8').split(",")
        # # logging.info('Current residents: ' + str(hotel_residents))
        # if hotel_residents == "":
        #     client.save_state(DAPR_STORE_NAME, hotel_key, init_value) 
        
        # new_residents = hotel_residents + [str(user.id)]
        # client.save_state(DAPR_STORE_NAME, hotel_key, ",".join(new_residents)) 

    ## No DB and Dapr alternative
    hotel_residents = []

    return {"hotel_id": hotel_id, "user_id": user.id, "new_occupancy": len(hotel_residents)}



# app.run(host="0.0.0.0",port=appPort)
