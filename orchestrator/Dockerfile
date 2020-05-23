FROM videra/gopyzmq-alpine

WORKDIR /app
COPY . /app

RUN virtualenv --python=python3 venv
RUN source ./venv/bin/activate

ENTRYPOINT ["make", "run"]
