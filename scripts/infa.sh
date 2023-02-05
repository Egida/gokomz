#!/bin/bash

docker rm $(docker ps -a -q)
docker volume rm $(docker volume ls -q)

make build && make container && docker-compose up
