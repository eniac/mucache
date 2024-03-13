FROM golang:alpine

RUN mkdir /app
COPY ./cmd /app/cmd
COPY ./internal /app/internal
COPY ./pkg /app/pkg
COPY ./go.mod /app/
COPY ./go.sum /app/
WORKDIR /app

RUN go build -o /app/main /app/cmd/singleservice/main.go

CMD /app/main
