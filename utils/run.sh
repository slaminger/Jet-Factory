#!/bin/bash

#Ensure nothing happens outside the parent directory
cd "$(dirname "$0")"
cd ..

SCRIPT_DIRECTORY=$(pwd)

export DEBIAN_FRONTEND=noninteractive

./jetfactory
