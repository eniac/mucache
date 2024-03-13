#!/usr/bin/env python3

import argparse
import copy
import logging
import json
import os
import requests
import random
import string
import sys

import concurrent.futures

from progress.bar import Bar

BAR_GRANULARITY = 100
BATCH_SIZE = 1000

## Taken from https://stackoverflow.com/questions/10272879/how-do-i-import-a-python-script-from-a-sibling-directory
sys.path.append(
    os.path.normpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', '..', 'tests', 'social')))
import utility


def get_random_string(length: int) -> string:
    # choose from all lowercase letter
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for i in range(length))
    # print("Random string of length", length, "is:", result_str)
    return result_str


def get_random_user_id(number_of_users: int) -> string:
    user_id_int = random.randrange(number_of_users)
    return f'user_{user_id_int}'

SERVICES = ["frontend", "product_catalog", "currency"]


def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--products_file",
                        help="the path that contains the social_graph_file",
                        default="./experiments/boutique/data/products.json")
    parser.add_argument("--currencies_file",
                        help="the path that contains the social_graph_file",
                        default="./experiments/boutique/data/currency_conversion.json")
    parser.add_argument("--analysis_file",
                        help="the path that contains the analysis of the social_graph",
                        default="./experiments/boutique/data/analysis.txt")
    parser.add_argument("--users",
                        help="number of users",
                        type=int,
                        default=10000)
    parser.add_argument("--products",
                        help="number of products",
                        type=int,
                        default=100)
    parser.add_argument("--home_catalog_size",
                        help="number of products in the home catalog",
                        type=int,
                        default=10)
    parser.add_argument("--product_size",
                        help="size of each product object",
                        type=int,
                        default=1000)
    # parser.add_argument("--number_posts_per_user",
    #                     help="the number of posts to create per user",
    #                     type=int,
    #                     default=10)
    for service in SERVICES:
        parser.add_argument(f'--{service}',
                            help=f'ip of the {service} service',
                            type=str,
                            required=True)
    args = parser.parse_args()
    return args


def normalize_product_ids_and_store_analysis(raw_products, number_products, users, home_catalog_size, product_size, analysis_file):
    num_raw_products = len(raw_products['products'])
    products = []

    # print(raw_products['products'])
    image = get_random_string(product_size)
    ## Normalize ids to start from 0
    for i in range(number_products):
        # print(i)
        products.append(copy.deepcopy(raw_products['products'][i % num_raw_products]))
        products[i]['id'] = f'p{i}'
        products[i]['picture'] = image

    with open(analysis_file, "w") as f:
        f.write(f'Users: {users}\n')
        f.write(f'Products: {number_products}\n')
        f.write(f'Home catalog size: {home_catalog_size}\n')
    
    return products

def add_product_batch(product_batch):
    global IPS
    url = utility.compose_url(IPS, "product_catalog", 'add_products')
    data = {
        'products': product_batch
    }
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

def populate_currencies(currencies):
    global IPS
    url = utility.compose_url(IPS, "currency", "init_currencies")
    currency_list = []
    for curr_id, curr_val in currencies.items():
        currency = {
            'currencyCode': curr_id,
            'rate': curr_val,
        }
        currency_list.append(currency)
    data = {
        'currencies': currency_list
    }
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

def populate_everything(num_of_users, products, currencies):
    global IPS
    ## Perform a heartbeat
    utility.heartbeat(IPS, "frontend")

    populate_currencies(currencies)
    
    bar = Bar('Populating products', max=BAR_GRANULARITY)
    counter = 0
    total_products = len(products)
    product_batch = []
    
    for products in products:
        counter += 1
        product_batch.append(products)
        if counter % (total_products // BAR_GRANULARITY) == 0:
            add_product_batch(product_batch)
            product_batch = []
            bar.next()
    
    ## Final batch
    add_product_batch(product_batch)
    bar.finish()




def main(args):
    global IPS
    IPS = utility.general_ips_from_args(args)
    products_file = args.products_file
    currencies_file = args.currencies_file
    analysis_file = args.analysis_file
    users = args.users
    number_products = args.products
    home_catalog_size = args.home_catalog_size
    product_size = args.product_size

    print(f'Parsing products file: {products_file}... ', end="")
    with open(products_file, "r") as f:
        raw_products = json.load(f)
    print(f'DONE!')

    print(f'Parsing currencies file: {currencies_file}... ', end="")
    with open(currencies_file, "r") as f:
        raw_currencies = json.load(f)
    print(f'DONE!')

    print(f'Producing analysis file: {analysis_file}... ', end="")
    ## Normalize the ids of movies
    products = normalize_product_ids_and_store_analysis(raw_products, number_products, users, home_catalog_size, product_size, analysis_file)
    print(f'DONE!')

    print(f'Populating everything... ')
    populate_everything(users, products, raw_currencies)
    print(f'DONE!')



if __name__ == '__main__':
    args = parse_arguments()
    main(args)
