FROM golang:alpine

#RUN apk add --no-cache --update libzmq git go gcc g++ pkgconfig musl-dev zeromq-dev
RUN apk add --no-cache --update gcc musl-dev zeromq-dev iftop

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
ENV EXPIRATION_TTL=$EXPIRATION_TTL

CMD if [ "$CM_ENABLED" = "true" ]; then \
        CGO_ENABLED=1 GOOS=linux go run -tags k8s,node${NODE_IDX},cm /app/cmd/${APP_NAMESPACE}/${APP_NAME}/main.go -logtostderr=true; \
    elif [ "$CM_ENABLED" = "upper" ]; then \
        CGO_ENABLED=1 GOOS=linux go run -tags k8s,node${NODE_IDX},upperbound /app/cmd/${APP_NAMESPACE}/${APP_NAME}/main.go -logtostderr=true; \
    else \
        CGO_ENABLED=1 GOOS=linux go run -tags k8s,node${NODE_IDX} /app/cmd/${APP_NAMESPACE}/${APP_NAME}/main.go -logtostderr=true; \
    fi
