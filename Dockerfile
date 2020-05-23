FROM videra/gopyzmq-alpine

WORKDIR /app
COPY . /app

RUN virtualenv --python=python3 venv
RUN source ./venv/bin/activate

RUN ./venv/bin/pip install -r /app/executor/requirements.txt

WORKDIR /app/orchestrator

ENTRYPOINT ["make", "run"]
