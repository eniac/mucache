import argparse
import logging
import json
import os
import requests
import random
import string
import sys

import concurrent.futures

from tqdm import tqdm

BAR_GRANULARITY = 100

## Taken from https://stackoverflow.com/questions/10272879/how-do-i-import-a-python-script-from-a-sibling-directory
sys.path.append(
    os.path.normpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', '..', 'tests', 'social')))
import utility

IPS = None

SERVICES = ["frontend", "user"]

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--hotels_file",
                        help="the path that contains the hotel file",
                        default="./experiments/hotel/data/hotels.json")
    parser.add_argument("--analysis_file",
                        help="the path that contains the analysis file",
                        default="./experiments/hotel/data/analysis.txt")
    parser.add_argument("--num_of_users",
                        help="the number of users",
                        type=int,
                        default=100)
    parser.add_argument("--info_size",
                        help="the size of hotel info in bytes",
                        type=int,
                        default=1000)
    for service in SERVICES:
        parser.add_argument(f'--{service}',
                            help=f'ip of the {service} service',
                            type=str,
                            required=True)
    args = parser.parse_args()
    return args

def get_random_string(length: int) -> string:
    # choose from all lowercase letter
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for i in range(length))
    # print("Random string of length", length, "is:", result_str)
    return result_str


def add_hotel(hotel, info_size):
    global IPS
    url = utility.compose_url(IPS, "frontend", 'store_hotel')
    data = {
        "hotel_id": str(hotel['id']),
        "name": hotel['name'],
        "phone": hotel['phoneNumber'],
        "location": hotel['address']['city'],
        "rate": 100,
        ## Note: The capacity is 11 to always allow reservations to succeed.
        ##       This works because we always leave at most 10 reservations for each hotel (for getState predictability).
        "capacity": 11,
        "info": get_random_string(info_size),
    }
    # print(data)
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

def upload_user(user_index):
    url = utility.compose_url(IPS, 'user', 'register_user')
    data = {
        "username": f'username{user_index}', 
        "password": f'password{user_index}'
    }
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r_ok = r.json()['ok']
    assert(r_ok == True)

def aggregate_location_popularities(hotels):
    locations = {}
    for hotel in hotels:
        city = hotel['address']['city']
        if city in locations:
            locations[city] += 1
        else:
            locations[city] = 1
    
    number_of_hotels = len(hotels)
    for key in locations.keys():
        locations[key] = locations[key] / number_of_hotels
    return locations

def normalize_and_store_analysis(hotels, users, analysis_file):
    ## Normalize the ids to start from 0
    for i in range(len(hotels)):
        hotels[i]['id'] = str(i)

    locations = aggregate_location_popularities(hotels)
    with open(analysis_file, "w") as f:
        f.write(f'Users: {users}\n')
        f.write(f'Hotels: {len(hotels)}\n')
        f.write(f'Cities: {len(locations)}\n')
        ## TODO: Maybe add popularity of searches of each city? in a different way?
        for location, popularity in locations.items():
            f.write(f'{location}|{popularity}\n')
    
    return hotels

def populate_everything(num_of_users: int, info_size: int, hotels):
    global IPS
    ## Perform a heartbeat
    utility.heartbeat(IPS, "frontend")

    ## TODO: Add users
    for i in range(num_of_users):
        upload_user(i)
    
    with concurrent.futures.ProcessPoolExecutor(max_workers=8) as executor:
        print("Populating hotels...")
        pending_futures = []
        for hotel in hotels:
            f = executor.submit(add_hotel, hotel, info_size)
            pending_futures.append(f)
        for f in tqdm(pending_futures):
            f.result()


def main(args):
    global IPS
    IPS = utility.general_ips_from_args(args)
    # print(args)
    # print(IPS)
    hotels_file = args.hotels_file
    analysis_file = args.analysis_file

    print(f'Parsing hotels file: {hotels_file}... ', end="")
    with open(hotels_file, "r") as f:
        raw_hotels = json.load(f)
    print(f'DONE!')

    print(f'Producing analysis file: {analysis_file}... ', end="")
    hotels = normalize_and_store_analysis(raw_hotels, args.num_of_users, analysis_file)
    print(f'DONE!')


    print(f'Populating everything... ')
    populate_everything(args.num_of_users, args.info_size, hotels)
    print(f'DONE!')


if __name__ == '__main__':
    args = parse_arguments()
    main(args)
