#!/bin/sh

echo "Run migrations"
migrate -path ./deploy/migrations -database $DATABASE_URL up && echo "Success"