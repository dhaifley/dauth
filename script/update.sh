#!/bin/sh
docker stack rm dauth
sleep 2s
docker pull dhaifley/dauth:latest
sleep 2s
docker system prune --volumes -f
sleep 2s
docker stack deploy -c /home/dhaifley/dauth/docker-compose.yml dauth
