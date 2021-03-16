#!/usr/bin/env bash

GOROOT=${GOROOT:-`go env GOROOT`}
GO=${GOROOT}/bin/go
APP_ROOT=${PWD}
MODULE=github.com/screwyprof/skeleton

# This script MUST be run from the project's root.

# This script creates fake test files so that 'go test' generates correct coverage profiles

FAKE_TEST_FILE=fake_for_correct_coverage_test.go

${GO} list -f '{{.ImportPath}} {{.Name}}' ./... | # list all packages including path and package name (sometimes they differ)
    sed 's/^skeleton\///' |                        # trim `skeleton/` prefix
    grep -v '^docs' |                             # exclude auto-generated docs
    grep -v '^tools' |                            # exclude tools - we're not going to additionally test them
    grep -v 'main$' |                             # exclude main package
    grep -v '^tests/\?' |                         # exclude all folders under tests/
while read in; do
    import=`echo $in | cut -d " " -f 1`
    package=`echo $in | cut -d " " -f 2`
    file=${APP_ROOT}${import//$MODULE/}
    ls ${file} | grep '_test.go$' > /dev/null
    if [ $? -eq 1 ]; then
        echo "Generated $file/$FAKE_TEST_FILE"
        echo "// DO NOT WRITE TESTS INTO THIS AUTO-GENERATED FILE, CREATE ANOTHER TEST FILE WITH A MEANINGFUL NAME" > ${file}/${FAKE_TEST_FILE}
        echo "package ${package}" >> ${file}/${FAKE_TEST_FILE}
    fi
done