#!/usr/bin/env bash

#see https://github.com/ory/go-acc
touch ./coverage.tmp
echo 'mode: atomic' > coverage.txt
go list ./... | grep -v /cmd | grep -v /vendor | xargs -n1 -I{} sh -c 'go test -race -covermode=atomic -coverprofile=coverage.tmp -coverpkg $(go list ./... | grep -v /vendor | tr "\n" ",") {} && tail -n +2 coverage.tmp >> coverage.txt || exit 255' && rm coverage.tmp
go test ./tests/... -tags="integration acceptance" -race -covermode=atomic -coverprofile=coverage_e2e.tmp -coverpkg $(go list ./... | grep -v /vendor | tr "\n" ",") {} && tail -n +2 coverage_e2e.tmp >> coverage.txt || exit 255 && rm coverage_e2e.tmp && cp coverage.txt c.out