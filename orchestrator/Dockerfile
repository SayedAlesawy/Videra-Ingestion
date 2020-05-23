FROM golang:1.14-alpine
WORKDIR /app
COPY . /app

RUN apk update && apk add build-base musl-dev libzmq zeromq-dev python3 py3-virtualenv

RUN virtualenv --python=python3 venv
RUN source ./venv/bin/activate

ENTRYPOINT ["make", "run"]
