#!/bin/sh -ex

docker build -t tiborvass/buildkit-ninja . && docker build -f ./test/build.ninja --progress plain --no-cache test/
