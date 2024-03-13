#!/bin/bash -ex

APP=$1

if [ -z "$APP" ]; then
  echo "Error: No app name provided"
  exit 1
fi

# set app_port in ports, dapr_http_port = app_port + 500
declare -A ports
ports["post_storage"]="3000"
ports["home_timeline"]="3001"
ports["user_timeline"]="3002"
ports["social_graph"]="3003"
ports["compose_post"]="3004"

LOC=$(dirname $(realpath "$0"))
dapr run --app-id "$APP" --app-port ${ports[$APP]} --dapr-http-port $((${ports[$APP]} + 500)) go run "$LOC"/"$APP"/main.go
