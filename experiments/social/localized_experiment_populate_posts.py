#!/usr/bin/env python3

import argparse
import logging
import os
import requests
import random
import string
import sys

## Taken from https://stackoverflow.com/questions/10272879/how-do-i-import-a-python-script-from-a-sibling-directory
sys.path.append(os.path.normpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', '..', 'tests', 'social')))
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

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("post_ids_file", 
                        help="the path to store post_ids in")
    parser.add_argument("--users", 
                        help="the number of users to create profiles for and posts",
                        type=int,
                        default=100)
    parser.add_argument("--posts", 
                        help="the number of posts to create",
                        type=int,
                        default=1000)
    args = parser.parse_args()
    return args

## TODO: Do we want some arguments (like distribution, users, etc)
def populate_posts(port: int, users: int, number_posts: int, post_ids_file: string):
    APP = "post_storage"
    PORT = port

    post_ids = []
    for p in range(number_posts):
        ## Check that a post that is written can be retrieved
        post_text = get_random_string(255)
        post_user = get_random_user_id(users)

        ## Write a post to post_storage
        url=f'http://localhost:{PORT}/v1.0/invoke/{APP}/method/store_post'
        data = {'text': post_text,
                'creator_id': post_user}
        r = requests.post(url, json=data)
        logging.debug(f'Store_post response code: {r.status_code}')
        assert(r.status_code == 200)
        response = r.json()
        post_id = response['post_id']
        logging.debug(f'Store_post response:\n{response}')

        ## Output post_id to save it somewhere to be able to retrieve it
        ## TODO: Store that in a file directly
        post_ids.append(post_id)
    
    post_id_lines = [f'{post_id}\n' for post_id in post_ids]
    post_id_file_data = "".join(post_id_lines)
    with open(post_ids_file, "w") as f:
        f.write(post_id_file_data)


def main(args):
    post_storage_port = args.post_storage_port
    number_of_users = args.users
    number_of_posts = args.posts
    post_ids_file = args.post_ids_file
    
    populate_posts(post_storage_port, number_of_users, number_of_posts, post_ids_file)

if __name__ == '__main__':
    args = parse_arguments()
    main(args)
