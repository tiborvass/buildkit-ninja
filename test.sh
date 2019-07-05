#!/bin/sh -x

docker build -t tiborvass/buildkit-ninja . || exit 1
start=$(date +%s)
docker build -f ./test/build.ninja --progress plain --no-cache test/
journalctl --since @$start -xu docker | grep totodebug
