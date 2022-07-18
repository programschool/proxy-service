#!/usr/bin/env bash

kill $(ps aux | grep lsof.sh | grep -v grep | awk '{print $2}')

docker-compose -f docker-compose.yaml up -d

bash monitor.sh
