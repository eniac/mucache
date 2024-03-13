#!/usr/bin/env python3

import argparse
import requests
import utility


def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--post_storage",
                        help="ip of the post_storage service",
                        type=str,
                        required=True)
    parser.add_argument("--home_timeline",
                        help="the dapr port of the home_timeline service",
                        type=str,
                        required=True)
    parser.add_argument("--user_timeline",
                        help="the dapr port of the user_timeline service",
                        type=str,
                        required=True)
    parser.add_argument("--social_graph",
                        help="ip of the social_graph service",
                        type=str,
                        required=True)
    parser.add_argument("--compose_post",
                        help="ip of the compose_post service",
                        type=str,
                        required=True)
    return parser.parse_args()


def main(args):
    app = "compose_post"
    ips = utility.populate_ips_from_args(args)

    ## Perform a heartbeat
    utility.heartbeat(ips, app)

    user1 = 'User1'
    user2 = 'User2'

    ## Add users
    url = utility.compose_url(ips, 'social_graph', 'insert_user')
    data = {'user_id': user1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    data = {'user_id': user2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    ## user1 follows user2
    url = utility.compose_url(ips, 'social_graph', 'follow')
    data = {'follower_id': user1, 'followee_id': user2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    ## user2 posts
    url = utility.compose_url(ips, app, 'compose_post')
    data = {'creator_id': user2, 'text': 'Hello World!'}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    ## Check post exists in user2's user timeline
    url = utility.compose_url(ips, 'user_timeline', 'ro_read_user_timeline')
    data = {'user_id': user2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    posts = r.json()['posts']
    assert (len(posts) == 1)
    post = posts[0]
    assert (post['creator_id'] == user2)
    assert (post['text'] == 'Hello World!')

    post_id = post['post_id']

    ## Check post exists in user1's home timeline
    url = utility.compose_url(ips, 'home_timeline', 'ro_read_home_timeline')
    data = {'user_id': user1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    posts = r.json()['posts']
    assert (len(posts) == 1)
    post = posts[0]
    assert (post['creator_id'] == user2)
    assert (post['text'] == 'Hello World!')
    assert (post['post_id'] == post_id)

    ## Check post exists in post_storage
    url = utility.compose_url(ips, 'post_storage', 'ro_read_post')
    data = {'post_id': post_id}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    post = r.json()['post']
    assert (post['creator_id'] == user2)
    assert (post['text'] == 'Hello World!')

    print("Compose post test passed successfully!")


if __name__ == '__main__':
    args = parse_arguments()
    main(args)
