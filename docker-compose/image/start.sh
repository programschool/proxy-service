#!/usr/bin/env bash


docker-compose -f docker-compose.yaml up -d

cd internal
docker-compose -f docker-compose-internal.yaml up -d
