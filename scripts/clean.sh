#!/bin/bash -x

export MUCACHE_TOP=${MUCACHE_TOP:-$(git rev-parse --show-toplevel --show-superproject-working-tree)}

export application_namespace=${1?application name not given, e.g., social}

declare -A all_services
all_services["singleservice"]="service"
all_services["twoservices"]="caller callee"
all_services["chain"]="service1 service2 service3 service4 backend"
all_services["chain3"]="service1 service2 backend"
all_services["chain4"]="service1 service2 service3 backend"
all_services["star"]="frontend backend1 backend2 backend3 backend4"
all_services["fanin"]="frontend1 frontend2 frontend3 frontend4 backend"
all_services["loadcm"]="stub loader"

all_services["movie"]="cast_info compose_review frontend movie_id movie_info movie_reviews page plot review_storage unique_id user user_reviews"
all_services["hotel"]="frontend profile rate reservation search user"
all_services["social"]="post_storage home_timeline user_timeline social_graph compose_post"
all_services["boutique"]="cart checkout currency email frontend payment product_catalog recommendations shipping"

## Services
for app_name in ${all_services[$application_namespace]}; do
  app_name_no_underscores=${app_name//_/}
  APP_NAME_NO_UNDERSCORES="$app_name_no_underscores" \
    envsubst <"${MUCACHE_TOP}/deploy/app.yaml" | kubectl delete -f -
done

services=(${all_services[$application_namespace]})

## Cache Manager
for idx in "${!services[@]}"; do
  NODE_IDX=$((idx + 1)) envsubst <"${MUCACHE_TOP}/deploy/cm/cm.yaml" | kubectl delete -f -
done
