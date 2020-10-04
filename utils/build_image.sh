#!/bin/bash

docker image build -t alizkan/jet-factory:latest "$(dirname "$(dirname "$(readlink -f "$0")")")"
