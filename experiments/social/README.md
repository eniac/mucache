## Running workload and experiments for the social network application

To run experiments for the social network application, use the following code:

```sh
## First make sure you have deployed the application using
./scripts/deploy.sh konstantinoskallas social
## where the first argument is the docker.io username that hosts
## the prebuild docker containers of the application service
## To build containers from source, just run (with your own docker.io username):
## `./scripts/build_and_push.sh konstantinoskallas social`

## Populate the database with a relevant social network
python3 experiments/social/populate.py

## Run an experiment
./wrk2/wrk -t2 -c4 -d20 -R50 --latency -s experiments/social/mixed_workload.lua http://localhost:8084/compose_post

## At the end of the experiment you can clean the application using:
./scripts/clean.sh social

```

__Q: The yaml files in Kubernetes have some form of loadbalancing: Do we actually need that?__
