#!/bin/bash

CONTAINER_NAME="jet-factory-ide"
WORKSPACE_DIRECTORY="//root/workspace"
GIT_REPO_TO_CLONE="https://github.com/Azkali/Jet-Factory.git"
IDE_PORT="9091"

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

    echo "Do you want to add the current Jet-Factory folder as a Volume to persist changes ?"
    select persist in yes no; do
        case $persist in
            yes)
            docker run -d -p $IDE_PORT:$IDE_PORT \
                        --volume //var/run/docker.sock:/var/run/docker.sock \
                        --volume $(dirname $(dirname $(readlink -fm $0))):/root/workspace \
                        --user root \
                        --workdir //root \
                        --name $CONTAINER_NAME \
                        ubuntu:19.10 \
                        bash -c "$START_COMMAND"
            ;;
            no)
            docker run -d -p $IDE_PORT:$IDE_PORT \
                        --volume //var/run/docker.sock:/var/run/docker.sock \
                        --user root \
                        --workdir //root \
                        --name $CONTAINER_NAME \
                        ubuntu:19.10 \
                        bash -c "$START_COMMAND"
            ;;
            *)
            echo "Invalid entry."
            exit
            ;;
        esac
        break
    done

    bashCommand "apt-get update"
    installDependency "bindfs"
    installDependency "fusermount"
    installDependency "git"
    installDependency "curl"
    installDependency "docker.io"
    installDependency "software-properties-common"
    bashCommand "add-apt-repository ppa:longsleep/golang-backports -y"
    bashCommand "apt-get update"
    installDependency "golang-go"

    #Install IDE - Code Server
    bashCommand "curl -fsSL https://code-server.dev/install.sh | sh"

    if [[ ${persist} == no ]]; then
        #Clone Repo Down
        bashCommand "git clone \"$GIT_REPO_TO_CLONE\" \"$WORKSPACE_DIRECTORY\""        
    fi
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
