#!/bin/bash
# This script launch docker to build android or linux

echo "Type the path to a folder you want to use as storage"
echo "can be a mounted external HDD"

if [ ! -d "$1" ]
then
    echo "Not a VALID Directory !"
    exit 0
fi

basepath="$(realpath "$1")"

docker image build -t alizkan/jet-factory:latest "$(dirname "$(dirname "$(readlink -fm "$0")")")"
docker run --privileged --rm -it -v "$basepath":/root/linux alizkan/jet-factory:latest
