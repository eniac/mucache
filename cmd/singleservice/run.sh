#!/bin/bash -ex

APP="singleservice"
LOC=$(dirname $(realpath "$0"))
dapr run --app-id "$APP" --app-port 8000 --dapr-http-port 8500 go run "$LOC"/"$APP"/main.go
