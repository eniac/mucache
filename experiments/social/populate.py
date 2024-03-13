#!/usr/bin/env python3

import argparse
import logging
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


def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--social_graph_file",
                        help="the path that contains the social_graph_file",
                        default="./experiments/social/socfb/socfb-Reed98.mtx")
    parser.add_argument("--analysis_file",
                        help="the path that contains the analysis of the social_graph",
                        default="./experiments/social/socfb/socfb-analysis.txt")
    parser.add_argument("--post_size",
                        help="the size of the post to create",
                        type=int,
                        default=100)
    parser.add_argument("--number_posts_per_user",
                        help="the number of posts to create per user",
                        type=int,
                        default=10)
    # parser.add_argument("--post_storage",
    #                     help="ip of the post_storage service",
    #                     type=str,
    #                     required=True)
    # parser.add_argument("--home_timeline",
    #                     help="the dapr port of the home_timeline service",
    #                     type=str,
    #                     required=True)
    # parser.add_argument("--user_timeline",
    #                     help="the dapr port of the user_timeline service",
    #                     type=str,
    #                     required=True)
    parser.add_argument("--social_graph",
                        help="ip of the social_graph service",
                        type=str,
                        required=True)
    parser.add_argument("--compose_post",
                        help="ip of the compose_post service",
                        type=str,
                        required=True)
    args = parser.parse_args()
    return args


## TODO: Do we want some arguments (like distribution, users, etc)
def populate_social_graph(ips, social_graph: utility.SocialGraph, 
                          total_followers: int, 
                          post_size: int,
                          number_posts_per_user: int):
    with concurrent.futures.ProcessPoolExecutor(max_workers=8) as executor:
        app = "social_graph"
        ## Perform a heartbeat
        utility.heartbeat(ips, app)

        user_ids = social_graph.get_nodes()

        ## Add follow relations
        bar = Bar('Processing follows', max=BAR_GRANULARITY)
        url = utility.compose_url(ips, app, 'follow_multi')
        users = 0
        results = []
        for user_id in user_ids:
            ## Get the incoming indices too
            followee_indexes = social_graph.get_outgoing_nodes(user_id)
            follower_indexes = social_graph.get_incoming_nodes(user_id)
            user_id_str: str = utility.create_user_id_from_int(user_id)

            follower_data = []
            for follower_index in follower_indexes:
                follower_id: str = utility.create_user_id_from_int(follower_index)
                follower_data.append(follower_id)

            followee_data = []
            for followee_index in followee_indexes:
                followee_id: str = utility.create_user_id_from_int(followee_index)
                followee_data.append(followee_id)

            data = {'user_id': user_id_str, 
                    'follower_ids': follower_data, 
                    'followee_ids': followee_data}
            # print("User:", user_id, "Followers:", len(follower_data), "followees:", len(followee_data))
            future = executor.submit(requests.post, url, json=data)
            results.append(future)
            users += 1
            if len(user_ids) > BAR_GRANULARITY and users % (len(user_ids) // BAR_GRANULARITY) == 0:
                bar.next()

        bar.finish()

        ## Add posts
        bar = Bar('Waiting follows', max=BAR_GRANULARITY)
        users = 0
        for future in concurrent.futures.as_completed(results):
            # r = requests.post(url, json=data)
            r = future.result()
            logging.debug(f'follow_multi response code: {r.status_code}')
            assert (r.status_code == 200)

            users += 1
            if len(user_ids) > BAR_GRANULARITY and users % (len(user_ids) // BAR_GRANULARITY) == 0:
                bar.next()

        bar.finish()

        # Add posts
        bar = Bar('Filling posts', max=BAR_GRANULARITY)
        users = 0
        url = utility.compose_url(ips, "compose_post", "compose_post_multi")
        
        results = []
        for user_id in user_ids:
            user_id_str: str = utility.create_user_id_from_int(user_id)

            data = {'creator_id': user_id_str, 
                    'text': get_random_string(post_size), 
                    'number': number_posts_per_user}
            # print("User:", user_id, "Followers:", len(follower_data), "followees:", len(followee_data))
            # r = requests.post(url, json=data)
            future = executor.submit(requests.post, url, json=data)
            results.append(future)

            users += 1
            if users % (len(user_ids) // BAR_GRANULARITY) == 0:
                bar.next()
        bar.finish()

        ## Add posts
        bar = Bar('Waiting posts', max=BAR_GRANULARITY)
        users = 0
        for future in concurrent.futures.as_completed(results):
            # r = requests.post(url, json=data)
            r = future.result()
            logging.debug(f'compose_post_multi response code: {r.status_code}')
            assert (r.status_code == 200)

            users += 1
            if len(user_ids) > BAR_GRANULARITY and users % (len(user_ids) // BAR_GRANULARITY) == 0:
                bar.next()
        bar.finish()


def main(args):
    ips = utility.populate_ips_from_args(args)
    post_size = args.post_size
    number_posts_per_user = args.number_posts_per_user
    social_graph_file = args.social_graph_file
    analysis_file = args.analysis_file

    print(f'Parsing social graph file: {social_graph_file}... ', end="")
    social_graph = utility.parse_social_graph(social_graph_file)
    print(f'DONE!')

    print(f'Analyzing social graph in: {analysis_file}... ', end="")
    total_followers = utility.analyze_social_graph(social_graph, analysis_file, post_size, number_posts_per_user)
    print(f'DONE!')

    print(f'Populating social graph... ')
    populate_social_graph(ips, social_graph, total_followers, post_size, number_posts_per_user)
    print(f'DONE!')

    ## TODO: Maybe we need to prepopulate with some compose_posts too


if __name__ == '__main__':
    args = parse_arguments()
    main(args)
