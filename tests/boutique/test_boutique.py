#!/usr/bin/env python3

import argparse
import os
import requests
import sys
import json

## Taken from https://stackoverflow.com/questions/10272879/how-do-i-import-a-python-script-from-a-sibling-directory
sys.path.append(
    os.path.normpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'social')))
import utility

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--product_catalog",
                        help="ip of the product_catalog service",
                        type=str,
                        required=True)
    parser.add_argument("--recommendations",
                        help="ip of the recommendations service",
                        type=str,
                        required=True)
    parser.add_argument("--cart",
                        help="ip of the cart service",
                        type=str,
                        required=True)
    parser.add_argument("--shipping",
                        help="ip of the shipping service",
                        type=str,
                        required=True)
    parser.add_argument("--payment",
                        help="ip of the payment service",
                        type=str,
                        required=True)
    parser.add_argument("--checkout",
                        help="ip of the checkout service",
                        type=str,
                        required=True)
    parser.add_argument("--email",
                        help="ip of the email service",
                        type=str,
                        required=True)
    return parser.parse_args()

def main(args):
    app = "product_catalog"
    ips = utility.general_ips_from_args(args)
    print(ips)

    ## Perform a heartbeat
    utility.heartbeat(ips, app)

    product_id1 = "OLJCESPC7Z"
    product_name1 = "Sunglasses"

    f = open("/users/pavlatos/mucache/cmd/boutique/products.json")
    json_data = json.load(f)
    product_ids = [product['id'] for product in json_data[:3]]
    
    url = utility.compose_url(ips, 'product_catalog', 'add_product')
    for product in json_data:
        data = {"product": product}
        r = requests.post(url, json=data)
        assert (r.status_code == 200)
 
    url = utility.compose_url(ips, 'product_catalog', 'ro_get_product')
    data = {'id': product_id1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    info1 = r.json()['product']
    assert(info1['name'] == product_name1)

    url = utility.compose_url(ips, 'product_catalog', 'ro_fetch_catalog')
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    url = utility.compose_url(ips, 'recommendations', 'ro_get_recommendations')
    data = {'user_id': "1", 'product_ids': product_ids}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    info1 = r.json()['product_ids']
    assert(info1[0] not in product_ids)

if __name__ == '__main__':
    args = parse_arguments()
    main(args)
