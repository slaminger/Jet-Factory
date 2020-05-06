#!/bin/env bash

# Setup variables
clean=false

# Android version specific variables
selection="$(echo ${@: -1} | tr '[:upper:]' '[:lower:]')"

(
	if [[ ${selection} == "icosa" ]]; then
		echo "Building Android Pie Icosa..."
	elif [[ ${selection} == "foster" ]]; then
		echo "Building Android Pie foster..."
	elif [[ ${selection} == "foster" ]]; then
		echo "Building Android Pie foster_tab..."
	else
		echo "$0: invalid distro option: $1"
		usage
		exit 1
	fi
)

Build() {
	git clone https://github.com/PabloZaiden/switchroot-android-build
	cd switchroot-android-build && mkdir -p ./android/lineage
	docker run --rm -ti -e ROM_NAME="${@: -1}" -v "$PWD"/android:/root/android pablozaiden/switchroot-android-build
}

# Parse arguments
# options=$(getopt -n $0 -o dhks --long docker,keep,staging,hekate,help -- "$@")
options=$(getopt -n $0 -o c --clean,help -- "$@")

# Check for errors in arguments or if no name was provided
if [[ $? != "0" ]] || [[ options =~ "${@: -1}" ]]; then usage; exit 1; fi

# Evaluate arguments
eval set -- "$options"
while true; do
    case "$1" in
	-c | --clean) clean=true; shift ;;
    ? | -h | --help) usage; exit 0 ;;
    -- ) shift; break ;;
    esac
done

Build