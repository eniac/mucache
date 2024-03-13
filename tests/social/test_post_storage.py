#!/usr/bin/env python3

import argparse
import logging
import requests

import utility


def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--post_storage",
                        help="ip of the post_storage service",
                        type=str,
                        required=True)
    return parser.parse_args()


def main(args):
    app = "post_storage"
    ips = utility.populate_ips_from_args(args)

    ## Perform a heartbeat
    utility.heartbeat(ips, app)

    ## Check that a post that is written can be retrieved
    post_text = 'Hello world, my first post'
    post_user = 'User0'

    ## Write a post to post_storage
    url = utility.compose_url(ips, app, 'store_post')
    data = {'text': post_text,
            'creator_id': post_user}
    r = requests.post(url, json=data)
    logging.debug(f'Store_post response code: {r.status_code}')
    assert (r.status_code == 200)
    response = r.json()
    post_id = response['post_id']
    logging.debug(f'Store_post response:\n{response}')

    ## Read post
    url = utility.compose_url(ips, app, 'read_post')
    data = {'post_id': post_id}
    r = requests.post(url, json=data)
    logging.debug(f'Read_post response code: {r.status_code}')
    assert (r.status_code == 200)

    response = r.json()
    logging.debug(f'Read_post response:\n{response}')
    returned_post_id = response['post']['post_id']
    returned_post_text = response['post']['text']
    returned_post_user = response['post']['creator_id']
    # print(response)
    assert (post_text == returned_post_text)
    assert (post_user == returned_post_user)
    assert (post_id == returned_post_id)

    print("Post storage test passed successfully!")


if __name__ == '__main__':
    args = parse_arguments()
    main(args)
