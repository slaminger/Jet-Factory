#!/bin/bash

CONTAINER_NAME="jet-factory-ide"
WORKSPACE_DIRECTORY="//root/workspace"
GIT_REPO_TO_CLONE="https://github.com/Azkali/Jet-Factory.git"
IDE_PORT="9090"

START_COMMAND="code-server --auth none --bind-addr 0.0.0.0:$IDE_PORT \"$WORKSPACE_DIRECTORY\" || tail -f /dev/null"

function bashCommand()
{
    COMMAND="$1"
    docker exec -it $CONTAINER_NAME bash -c "$COMMAND"
}

function installDependency()
{
    DEPENDENCY="$1"
    bashCommand "command -v $DEPENDENCY || apt-get install -y $DEPENDENCY"
}

function doSetup()
{
    docker stop $CONTAINER_NAME
    docker rm $CONTAINER_NAME

    docker run -d -p $IDE_PORT:$IDE_PORT \
                --volume //var/run/docker.sock:/var/run/docker.sock \
                --user root \
                --workdir //root \
                --name $CONTAINER_NAME \
                debian:latest \
                bash -c "$START_COMMAND"

    bashCommand "apt-get update"
    installDependency "bindfs"
    installDependency "fusermount"
    installDependency "git"
    installDependency "curl"
    installDependency "docker.io"

    #Install IDE - Code Server
    bashCommand "curl -fsSL https://code-server.dev/install.sh | sh"

    #Clone Repo Down
    bashCommand "git clone \"$GIT_REPO_TO_CLONE\" \"$WORKSPACE_DIRECTORY\""

    #Start IDE
    bashCommand "code-server --install-extension ms-azuretools.vscode-docker"
    bashCommand "code-server --install-extension ms-vscode.go"

    docker restart $CONTAINER_NAME
}

if [ ! "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    #Setup the Code Environment
    doSetup
else
    #Start the Code Environment (If it already exists)
    docker start $CONTAINER_NAME
fi
