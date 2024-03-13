FROM golang:alpine

#RUN apk add --no-cache --update libzmq git go gcc g++ pkgconfig musl-dev zeromq-dev
RUN apk add --no-cache --update gcc musl-dev zeromq-dev

RUN mkdir /app
COPY ./cmd /app/cmd
COPY ./internal /app/internal
COPY ./pkg /app/pkg
COPY ./go.mod /app/
COPY ./go.sum /app/
WORKDIR /app

ENV APP_NAME=$APP_NAME
ENV APP_NAMESPACE=$APP_NAMESPACE
ENV NODE_IDX=$NODE_IDX
ENV CM_ENABLED=$CM_ENABLED
ENV SHARD_IDX=$SHARD_IDX
ENV SHARD_COUNT=$SHARD_COUNT
ENV APP_RAW_NAME=$APP_RAW_NAME
ENV APP_RAW_NAME_NO_UNDERSCORES=$APP_RAW_NAME_NO_UNDERSCORES

CMD if [ "$CM_ENABLED" = "true" ]; then \
        CGO_ENABLED=1 GOOS=linux go run -tags k8s,node${NODE_IDX},cm,shard /app/cmd/${APP_NAMESPACE}/${APP_RAW_NAME}/main.go -logtostderr=true; \
    else \
        CGO_ENABLED=1 GOOS=linux go run -tags k8s,node${NODE_IDX},shard /app/cmd/${APP_NAMESPACE}/${APP_RAW_NAME}/main.go -logtostderr=true; \
    fi
