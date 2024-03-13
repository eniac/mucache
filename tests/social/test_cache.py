#!/usr/bin/env python3

import argparse
import requests
import utility


def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--compose_post",
                        help="ip of the compose storage service",
                        type=str,
                        required=True)
    parser.add_argument("--post_storage",
                        help="ip of the post storage service",
                        type=str,
                        required=True)
    parser.add_argument("--user_timeline",
                        help="ip of the user timeline service",
                        type=str,
                        required=True)
    return parser.parse_args()


def main(args):
    ips = utility.populate_ips_from_args(args)

    url = utility.compose_url(ips, "compose_post", 'compose_post')
    data = {'creator_id': "User1", 'text': 'Hello World!'}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    # cache miss
    url = utility.compose_url(ips, 'user_timeline', 'ro_read_user_timeline')
    data = {'user_id': "User1"}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    posts = r.json()['posts']
    assert (len(posts) == 1)

    # cache hit
    url = utility.compose_url(ips, 'user_timeline', 'ro_read_user_timeline')
    data = {'user_id': "User1"}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    posts = r.json()['posts']
    assert (len(posts) == 1)

    url = utility.compose_url(ips, "compose_post", 'compose_post')
    data = {'creator_id': "User1", 'text': 'Hello World!'}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    # cache miss
    url = utility.compose_url(ips, 'user_timeline', 'ro_read_user_timeline')
    data = {'user_id': "User1"}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    posts = r.json()['posts']
    assert (len(posts) == 2)


if __name__ == '__main__':
    args = parse_arguments()
    main(args)
