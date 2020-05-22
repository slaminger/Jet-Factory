#!/bin/bash
# This script launch docker to build android or linux

echo "Type the path to a folder you want to use as storage"
echo "can be a mounted external HDD"

echo "Press enter..."
read basepath

if [ ! -d $basepath ] 
then
    echo "Directory $basepath DOES NOT exists. Exiting." 
    exit 1
fi

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

docker image build -t azkali/jet-factory:1.0.0 .

select distro in "$ARCH1" "$ARCH2" "$ARCH3" "$FEDORA" "$GENTOO" "$UBUNTU" "$LINEAGE" "$ICOSA" "$FOSTER" "$FOSTER_TAB"
do
    case $distro in
        $ARCH1)
            echo -e "building $ARCH1"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$ARCH1" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        $ARCH2)
            echo -e "building $ARCH2"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$ARCH2" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        $ARCH3)
            echo -e "building $ARCH3"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$ARCH3" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        $FEDORA)
            echo -e "building $FEDORA"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$FEDORA" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        $GENTOO)
            echo -e "building $GENTOO"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$GENTOO" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        $UBUNTU)
            echo -e "building $UBUNTU"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$UBUNTU" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        $LINEAGE)
            echo -e "building $ICOSA"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$ICOSA" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        $ICOSA)
            echo -e "building $ICOSA"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$ICOSA" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        $FOSTER)
            echo -e "building $FOSTER"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$FOSTER" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        $FOSTER_TAB)
            echo "building $FOSTER_TAB"
            docker run --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -it -e DISTRO="$FOSTER_TAB" -v "$basepath":/root/ -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0
            exit 0
        ;;
        *) 
            echo -e "\n ==> Enter a number between 1 and 10"
            exit 1
        ;;
    esac
done
