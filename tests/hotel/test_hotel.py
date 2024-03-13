#!/usr/bin/env python3

import argparse
import os
import requests
import sys

## Taken from https://stackoverflow.com/questions/10272879/how-do-i-import-a-python-script-from-a-sibling-directory
sys.path.append(
    os.path.normpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'social')))
import utility

SERVICES = ["frontend", "reservation", "user"]


def parse_arguments():
    parser = argparse.ArgumentParser()
    for service in SERVICES:
        parser.add_argument(f'--{service}',
                            help=f'ip of the {service} service',
                            type=str,
                            required=True)

    return parser.parse_args()

def make_hotel_dict(i):
    h_dict = {
        "hotel_id": str(i),
        "name": f'hotel{i}',
        "phone": f'phone:{i}',
        "location": 'Philadelphia',
        "rate": i * 100,
        "capacity": 1,
    }
    return h_dict

def main(args):
    app = "frontend"
    ips = utility.general_ips_from_args(args)

    ## Perform a heartbeat
    utility.heartbeat(ips, app)

    num_hotels = 4
    hotels = {str(i): make_hotel_dict(i) for i in range(num_hotels)}

    ## Add users
    username = "user"
    password = "pass"
    ## Register with users and login
    url = utility.compose_url(ips, 'user', 'register_user')
    data = {"username": username, "password": password}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r_ok = r.json()['ok']
    assert(r_ok == True)
    url = utility.compose_url(ips, 'user', 'login')
    data = {"username": username, "password": password}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r_token = r.json()['token']
    assert(r_token == "OK")

    ## Add hotel profiles
    for i in range(num_hotels):
        url = utility.compose_url(ips, 'frontend', 'store_hotel')
        data = hotels[str(i)]
        r = requests.post(url, json=data)
        assert (r.status_code == 200)
    
    
    url = utility.compose_url(ips, 'frontend', 'ro_search_hotels')
    data = {
        "in_date": "2023-04-01", 
        "out_date": "2023-04-02",
        "location": "Philadelphia"}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    print(r.json())
    r_hotels = r.json()['profiles']
    ## This fails if redis is not cleared
    assert(len(r_hotels) == len(hotels))
    for r_hotel in r_hotels:
        assert(hotels[r_hotel['hotel_id']]['hotel_id'] == r_hotel['hotel_id'])
        assert(hotels[r_hotel['hotel_id']]['name'] == r_hotel['name'])
        assert(hotels[r_hotel['hotel_id']]['phone'] == r_hotel['phone'])

    url = utility.compose_url(ips, 'frontend', 'reservation')
    data = {
        "in_date": "2023-04-01", 
        "out_date": "2023-04-02",
        "hotel_id": "0",
        "username": username,
        "password": password,
        "rooms": 1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    print(r.json())
    assert(r.json()["success"] == True)

    url = utility.compose_url(ips, 'frontend', 'reservation')
    data = {
        "in_date": "2023-04-01", 
        "out_date": "2023-04-02",
        "hotel_id": "0",
        "username": username,
        "password": password,
        "rooms": 1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    print(r.json())
    assert(r.json()["success"] == False)

    print("Hotel application test passed successfully!")


if __name__ == '__main__':
    args = parse_arguments()
    main(args)
