#!/bin/bash
# This script launch docker to build android or linux

echo "Type the path to a folder you want to use as storage"
echo "can be a mounted external HDD"

if [ ! -d $1 ]
then
    echo "Not a VALID Directory !"
    exit 0
fi

basepath="$PWD/$1"

docker image build -t alizkan/jet-factory:1.0.0 .
docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --rm -it -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0 -hekate
