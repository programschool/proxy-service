#!/usr/bin/env bash

p=$(ps aux | grep ./entry-proxy | grep -v grep | awk '{print $2}')

echo -e "pid=$p"

while [ true ];
do  
  sleep 2
  count=$(lsof -p $p | wc -l)
  if (( $count>3000 ))
  then
    docker restart entry-proxy
    p=$(ps aux | grep ./entry-proxy | grep -v grep | awk '{print $2}')
    echo -e "pid=$p"
    echo restart
  fi
done
