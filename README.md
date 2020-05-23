# Videra-Ingestion

[![Build Status:](https://github.com/SayedAlesawy/Videra-Ingestion/workflows/Build/badge.svg)](https://github.com/SayedAlesawy/Videra-Ingestion/actions)

The Ingestion Module Implementation for the Videra Video Indexer Engine.

## Setup

## Executor[Python]

```shell
virtualenv --python=python3 venv
source ./venv/bin/activate
python main.py
```

## Orchestrator

```shell
make build
make run
```
## Run Container
The following will build the docker image, establish a network and like the container on to that network.
```
docker-compose up ingestion-engine
```