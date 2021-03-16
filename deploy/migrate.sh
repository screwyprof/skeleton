#!/bin/sh

until nc -z database 5432;
do
  echo "Waiting for Postgres..." &&
  sleep 1;
done && migrate "$@"
