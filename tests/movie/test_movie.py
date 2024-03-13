#!/usr/bin/env python3

import argparse
import os
import requests
import sys

## Taken from https://stackoverflow.com/questions/10272879/how-do-i-import-a-python-script-from-a-sibling-directory
sys.path.append(
    os.path.normpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'social')))
import utility


def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--cast_info",
                        help="ip of the cast_info service",
                        type=str,
                        required=True)
    parser.add_argument("--compose_review",
                        help="ip of the compose_review service",
                        type=str,
                        required=True)
    parser.add_argument("--frontend",
                        help="ip of the frontend service",
                        type=str,
                        required=True)
    parser.add_argument("--movie_id",
                        help="ip of the movie_id service",
                        type=str,
                        required=True)
    parser.add_argument("--movie_info",
                        help="ip of the movie_info service",
                        type=str,
                        required=True)
    parser.add_argument("--movie_reviews",
                        help="ip of the movie_reviews service",
                        type=str,
                        required=True)
    parser.add_argument("--page",
                        help="ip of the page service",
                        type=str,
                        required=True)
    parser.add_argument("--plot",
                        help="ip of the plot service",
                        type=str,
                        required=True)

    parser.add_argument("--review_storage",
                        help="ip of the review_storage service",
                        type=str,
                        required=True)
    parser.add_argument("--unique_id",
                        help="ip of the unique_id service",
                        type=str,
                        required=True)
    parser.add_argument("--user",
                        help="ip of the user service",
                        type=str,
                        required=True)
    parser.add_argument("--user_reviews",
                        help="ip of the user_reviews service",
                        type=str,
                        required=True)
    return parser.parse_args()

def make_review_dict(i):
    r_dict = {
        "review_id": f'review{i}',
        "user_id":  'user',
        "req_id":  f'request{i}',
        "text": f'review text{i}',
        "movie_id": 'movie1',
        "rating": 10,
        "timestamp": 100
    }
    return r_dict

def main(args):
    app = "movie_info"
    ips = utility.general_ips_from_args(args)

    ## Perform a heartbeat
    utility.heartbeat(ips, app)

    cast_dict = {
        "c1": {"cast_id": "c1", "name": "cast1", "info": "cast 1 awards"},
        "c2": {"cast_id": "c2", "name": "cast2", "info": "cast 2 awards"},
        "c3": {"cast_id": "c3", "name": "cast3", "info": "cast 3 awards"},
    }

    movie_id1 = "movie1"
    movie_title1 = "jurassic1"
    movie_info1 = "This is a movie about dinosaurs"
    cast1 = ["c1", "c2"]

    movie_id2 = "movie2"
    movie_title2 = "planet2"
    movie_info2 = "This is a movie about planets"
    cast2 = ["c2", "c3"]

    username = "user"
    password = "pass"

    plot_id1 = "plot1"
    plot1 = "plot1: todo"

    plot_id2 = "plot2"
    plot2 = "plot2: todo"

    reviews_dict = { f'review{i}': make_review_dict(i) for i in range(3)}

    ## Add and get movie infos
    url = utility.compose_url(ips, 'movie_info', 'store_movie_info')
    data = {'movie_id': movie_id1, 'movie_info': movie_info1, 'cast_ids': cast1, "plot_id": plot_id1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    data = {'movie_id': movie_id2, 'movie_info': movie_info2, 'cast_ids': cast2, "plot_id": plot_id2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    url = utility.compose_url(ips, 'movie_info', 'ro_read_movie_info')
    data = {'movie_id': movie_id1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    info1 = r.json()['movie_info']
    print(info1, r.json())
    assert(info1['info'] == movie_info1)
    assert(info1['cast_ids'] == cast1)

    ## Register movie id with title
    url = utility.compose_url(ips, 'movie_id', 'register_movie_id')
    data = {'title': movie_title1, 'movie_id': movie_id1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    data = {'title': movie_title2, 'movie_id': movie_id2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    url = utility.compose_url(ips, 'movie_id', 'ro_get_movie_id')
    data = {'title': movie_title1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    info1 = r.json()['movie_id']
    print(info1, r.json())
    assert(info1 == movie_id1)

    ## Add and get cast infos
    url = utility.compose_url(ips, 'cast_info', 'store_cast_info')
    for cast_member in cast_dict.values():
        data = cast_member
        r = requests.post(url, json=data)
        assert (r.status_code == 200)

    url = utility.compose_url(ips, 'cast_info', 'ro_read_cast_infos')
    cast_ids = list(cast_dict.keys())
    data = {"cast_ids": cast_ids}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r_cast_infos = r.json()['cast_infos']
    print(r_cast_infos)
    assert(len(r_cast_infos) == len(list(cast_dict.keys())))
    for item in r_cast_infos:
        r_id = item["cast_id"]
        assert(item == cast_dict[r_id])

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

    ## Add and get plots
    url = utility.compose_url(ips, 'plot', 'write_plot')
    data = {'plot_id': plot_id1, 'plot': plot1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    data = {'plot_id': plot_id2, 'plot': plot2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)

    url = utility.compose_url(ips, 'plot', 'ro_read_plot')
    data = {'plot_id': plot_id1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r_plot = r.json()['plot']
    print(r.json())
    assert(r_plot == plot1)


    ## TODO: Add a test that performs compose in the frontend to compose reviews
    url = utility.compose_url(ips, 'frontend', 'compose')
    for review in reviews_dict.values():
        data = {
            "username": username,
            "password": password,
            "title": movie_title1,
            "rating": 10,
            "text": "Great movie about dinosaurs, highly recommend"
        }
        r = requests.post(url, json=data)
        assert (r.status_code == 200)
    


    ## Test that the page contains everything it needs to contain for a specific movie
    url = utility.compose_url(ips, 'page', 'ro_read_page')
    data = {'movie_id': movie_id1}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    info1 = r.json()['page']
    print(info1, r.json())
    assert(info1['movie_info']['info'] == movie_info1)
    assert(info1['plot'] == plot1)
    assert(len(info1['cast_infos']) == 2)
    assert(len(info1['reviews']) == len(reviews_dict))
    data = {'movie_id': movie_id2}
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    info2 = r.json()['page']
    assert(info2['movie_info']['info'] == movie_info2)
    assert(info2['plot'] == plot2)
    assert(len(info2['cast_infos']) == 2)
    assert(len(info2['reviews']) == 0)




    print("Movie application test passed successfully!")


if __name__ == '__main__':
    args = parse_arguments()
    main(args)
