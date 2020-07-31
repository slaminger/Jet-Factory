#!/usr/bin/bash
# ENTRYPOINT.SH : Manages options, and launches the sub scripts.

usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
	echo " -d, --docker          	Build with Docker"
	echo " -f, --force             	Download setup files anyway"
	echo " --hekate                 Build for Hekate"
    echo " -h, --help               Show this help text"
}

# Parse arguments
options=$(getopt -n $0 -o dfh --long docker,force,hekate,help -- "$@")

# Check for errors in arguments or if no name was provided
if [[ $? != "0" ]]; then usage; exit 1; fi

# Evaluate arguments
eval set -- "$options"
while true; do
    case "$1" in
	-a | --android) ROM_NAME=${OPTARG}; shift ;;
	-f | --force) force=true; shift ;;
	-d | --docker) docker=true; shift ;;
	--hekate)  hekate=true; shift ;;
    ? | -h | --help) usage; exit 0 ;;
    -- ) shift; break ;;
    esac
done

if [ ${ROM_NAME} != "" ]; then
	mkdir -p ./android/lineage
	docker run --rm -ti -e ROM_NAME=${ROM_NAME} -v "$PWD"/android:/root/android pablozaiden/switchroot-android-build:latest
	exit 0
fi

if [[ ${docker} == true ]]; then
	echo "Running container..."
	docker run --privileged --rm -it -v ${basepath}:/buildroot -v ${cwd}:/jetfactory -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:latest /jetfactory/src/entrypoint.sh "$(echo "$options" | sed -E 's/-(d|-docker)//g')"
	exit 0
fi

./net/dl_file.sh

if [ ${check} = "true" ]; then
	./net/checksum.sh
fi

./fs/extract_rootfs.sh

[[ hekate == "true" ]] && ./fs/hekate.sh