import argparse
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

## Taken from https://stackoverflow.com/questions/10272879/how-do-i-import-a-python-script-from-a-sibling-directory
sys.path.append(
    os.path.normpath(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', '..', 'tests', 'social')))
import utility

IPS = None


def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--movies_file",
                        help="the path that contains the movie file",
                        default="./experiments/movie/data/movies_1_500.json")
    parser.add_argument("--cast_file",
                        help="the path that contains the cast",
                        default="./experiments/movie/data/casts_1_500.json")
    parser.add_argument("--analysis_file",
                        help="the path that contains the analysis file",
                        default="./experiments/movie/data/analysis.txt")
    parser.add_argument("--num_of_users",
                        help="the number of users",
                        type=int,
                        default=100)
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
    args = parser.parse_args()
    return args


def upload_movie_info(movie):
    global IPS
    url = utility.compose_url(IPS, "movie_info", 'store_movie_info')
    data = {
        'movie_id': str(movie['id']),
        'movie_info': movie['title'],
        'cast_ids': [str(cast['id']) for cast in movie['cast']],
        'plot_id': str(movie['id'])
    }
    # print(data)
    r = requests.post(url, json=data)
    assert (r.status_code == 200)


def register_movie_id(movie):
    global IPS
    url = utility.compose_url(IPS, 'movie_id', 'register_movie_id')
    data = {
        'title': movie['title'],
        'movie_id': str(movie['id']),
    }
    # print(data)
    r = requests.post(url, json=data)
    assert (r.status_code == 200)


def write_plot(movie):
    global IPS
    url = utility.compose_url(IPS, 'plot', 'write_plot')
    data = {
        'plot_id': str(movie['id']),
        'plot': (movie['overview'] * 10),
    }
    # print(data)
    r = requests.post(url, json=data)
    assert (r.status_code == 200)


def upload_cast_info(cast):
    global IPS
    url = utility.compose_url(IPS, "cast_info", 'store_cast_info')
    data = {
        'cast_id': str(cast['id']),
        'name': cast['name'],
        'info': cast['biography']
    }
    # print(data['cast_id'])
    r = requests.post(url, json=data)
    assert (r.status_code == 200)


def upload_user(user_index):
    url = utility.compose_url(IPS, 'user', 'register_user')
    data = {
        "username": f'username{user_index}',
        "password": f'password{user_index}'
    }
    r = requests.post(url, json=data)
    assert (r.status_code == 200)
    r_ok = r.json()['ok']
    assert (r_ok == True)


def normalize_movie_ids_and_store_analysis(movies, users, casts, analysis_file):
    ## Normalize ids to start from 0
    for i in range(len(movies)):
        movies[i]['id'] = str(i)

    with open(analysis_file, "w") as f:
        f.write(f'Users: {users}\n')
        f.write(f'Movies: {len(movies)}\n')
        f.write(f'Cast members: {len(casts)}\n')
        f.write('Popularities|Titles:\n')
        for movie in movies:
            f.write(f'{movie["popularity"]}|{movie["title"]} \n')

    return movies


def populate_everything(num_of_users, movies, cast):
    global IPS
    ## Perform a heartbeat
    utility.heartbeat(IPS, "movie_info")

    ## TODO: Add users
    for i in range(num_of_users):
        upload_user(i)

    with concurrent.futures.ThreadPoolExecutor(max_workers=4) as executor:

        bar = Bar('Populating movies', max=BAR_GRANULARITY)
        counter = 0
        total_movies = len(movies)
        pending_futures = []
        for movie in movies:
            f = executor.submit(upload_movie_info, movie)
            pending_futures.append(f)
            f = executor.submit(register_movie_id, movie)
            pending_futures.append(f)
            f = executor.submit(write_plot, movie)
            pending_futures.append(f)
            counter += 1
            if counter % (total_movies // BAR_GRANULARITY) == 0:
                bar.next()
        bar.finish()

        bar = Bar('Waiting for movies', max=BAR_GRANULARITY)
        counter = 0
        total_futures = len(pending_futures)
        for f in pending_futures:
            f.result()
            counter += 1
            if counter % (total_futures // BAR_GRANULARITY) == 0:
                bar.next()
        bar.finish()

        bar = Bar('Populating cast', max=BAR_GRANULARITY)
        counter = 0
        pending_futures = []
        total_cast = len(cast)
        for cast_member in cast:
            f = executor.submit(upload_cast_info, cast_member)
            pending_futures.append(f)
            counter += 1
            if counter % (total_cast // BAR_GRANULARITY) == 0:
                bar.next()
        bar.finish()

        print("Pending cast:", len(pending_futures))

        bar = Bar('Waiting for cast', max=BAR_GRANULARITY)
        counter = 0
        total_futures = len(pending_futures)
        for f in pending_futures:
            try:
                f.result()
            except Exception as e:
                with open("error.log", "a") as ef:
                    ef.write(str(e))
                    ef.write("\n")
            counter += 1
            if counter % (total_futures // BAR_GRANULARITY) == 0:
                bar.next()
        bar.finish()


def main(args):
    global IPS
    IPS = utility.general_ips_from_args(args)
    # print(args)
    # print(IPS)
    movies_file = args.movies_file
    cast_file = args.cast_file
    analysis_file = args.analysis_file

    print(f'Parsing movies file: {movies_file}... ', end="")
    with open(movies_file, "r") as f:
        raw_movies = json.load(f)
    print(f'DONE!')

    print(f'Parsing cast file: {cast_file}... ', end="")
    with open(cast_file, "r") as f:
        cast = json.load(f)
    print(f'DONE!')

    print(f'Producing analysis file: {analysis_file}... ', end="")
    ## Normalize the ids of movies
    movies = normalize_movie_ids_and_store_analysis(raw_movies, args.num_of_users, cast, analysis_file)
    print(f'DONE!')

    print(f'Populating everything... ')
    populate_everything(args.num_of_users, movies, cast)
    print(f'DONE!')


if __name__ == '__main__':
    args = parse_arguments()
    main(args)
