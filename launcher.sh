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

ARCH1="arch"
ARCH2="blackarch"
ARCH3="arch-bang"
FEDORA="fedora"
GENTOO="gentoo"
UBUNTU="ubuntu"
LINEAGE="lineage (defaults to icosa)"
ICOSA="icosa"
FOSTER="foster"
FOSTER_TAB="foster_tab"

sudo rm -rf "$basepath"/*/disk "$basepath"/*/*.img 

select distro in "$ARCH1" "$ARCH2" "$ARCH3" "$FEDORA" "$GENTOO" "$UBUNTU" "$LINEAGE" "$ICOSA" "$FOSTER" "$FOSTER_TAB"
do   
    docker build -t alizkan/jet-factory:1.0.0 .
    case $distro in
        $ARCH1)
            echo -e "\nBuilding $ARCH1"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --rm -it -e DISTRO="$ARCH1" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        $ARCH2)
            echo -e "\nBuilding $ARCH2"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --rm -it -e DISTRO="$ARCH2" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        $ARCH3)
            echo -e "\nBuilding $ARCH3"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --rm -it -e DISTRO="$ARCH3" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        $FEDORA)
            echo -e "\nBuilding $FEDORA"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --rm -it -e DISTRO="$FEDORA" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        $GENTOO)
            echo -e "\nBuilding $GENTOO"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --rm -it -e DISTRO="$GENTOO" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        $UBUNTU)
            echo -e "\nBuilding $UBUNTU"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --rm -it -e DISTRO="$UBUNTU" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        $LINEAGE)
            echo -e "\nBuilding $ICOSA"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --privileged --rm -it -e DISTRO="$ICOSA" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        $ICOSA)
            echo -e "\nBuilding $ICOSA"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --privileged --rm -it -e DISTRO="$ICOSA" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        $FOSTER)
            echo -e "\nBuilding $FOSTER"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --privileged --rm -it -e DISTRO="$FOSTER" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        $FOSTER_TAB)
            echo "\nBuilding $FOSTER_TAB"
            docker run --name jet --privileged --cap-add=ALL --device=/dev/fuse --security-opt apparmor:unconfined --privileged --rm -it -e DISTRO="$FOSTER_TAB" -v "$basepath":/root/linux -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0
            exit 0
        ;;
        *) 
            echo -e "\n ==> Enter a number between 1 and 10"
            exit 1
        ;;
    esac
done
