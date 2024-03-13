#!/usr/bin/env python3

import argparse
import requests
import utility


def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--social_graph",
                        help="ip of the social_graph service",
                        type=str,
                        required=True)
    return parser.parse_args()


def main(args):
    app = "social_graph"
    ips = utility.populate_ips_from_args(args)

    ## Perform a heartbeat
    utility.heartbeat(ips, app)

    user1 = 'User1'
    user2 = 'User2'

    ## Add users
    url = utility.compose_url(ips, app, 'insert_user')
    data = {'user_id': user1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    data = {'user_id': user2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    ## user1 follows user2
    url = utility.compose_url(ips, app, 'follow')
    data = {'follower_id': user1, 'followee_id': user2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    ## check followers and followees
    url = utility.compose_url(ips, app, 'get_followers')
    ## user1 should have no followers
    data = {'user_id': user1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r = r.json()
    assert (r['followers'] == [])
    ## user2 should have user1 as a follower
    data = {'user_id': user2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r = r.json()
    assert (r['followers'] == [user1])

    url = utility.compose_url(ips, app, 'get_followees')
    ## user1 should have user2 as followee
    data = {'user_id': user1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r = r.json()
    assert (r['followees'] == [user2])
    ## user2 should have user1 as a follower
    data = {'user_id': user2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r = r.json()
    assert (r['followees'] == [])
    print("Social graph test passed successfully!")


if __name__ == '__main__':
    args = parse_arguments()
    main(args)
