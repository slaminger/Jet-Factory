#!/bin/bash

# ENTRYPOINT.SH : Manages options, and launches the sub scripts.

# Store the script directory
cwd=$(dirname "$(readlink -f "$0")")

# Read configs dir for config files
unset options i
while IFS= read -r -d $'\0' f; do
  options[i++]="$f"
done < <(find $(dirname ${cwd})/configs/ -maxdepth 1 -type f -print0 )

# Select DISTRO from configs files
select opt in "${options[@]}"; do
	source ${opt}
	export $(cut -d= -f1 ${opt})
	img="${URL##*/}"
break;
done

# Script usage
usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
	echo " -h, --hekate                 Build for Hekate"
    echo " -u, --usage               Show this help text"
}

# Parse arguments
options=$(getopt -n $0 -o hu --long hekate,usage -- "$@")

# Check for errors in arguments or if no name was provided
if [[ $? != "0" ]]; then usage; exit 1; fi

# Evaluate arguments
eval set -- "$options"
while true; do
    case "$1" in
	-h | --hekate)  hekate=true; shift ;;
    ? | -u | --usage) usage; exit 0 ;;
    -- ) shift; break ;;
    esac
done

# Get the last argument passed to the script to evaluate it later
out=$(realpath ${@:$#})

# The last argument must be a path pointing to a dir if not inside Docker
if [[ ! -f /.dockerenv ]] && [[ ! -d ${out} ]]; then
	usage
	exit 1
# If in docker, and not using --android flag, use the volume for linux
elif [[ -f /.dockerenv ]]; then
	out="/root/linux"
fi

cd ${out}

# Create the build directories if they don't exist
mkdir -p ${out}/{${NAME},downloadedFiles}

# Make scripts executable
chmod +x ${cwd}/{net,fs}/* $(dirname ${cwd})/configs/examples/*

# Download the URL
source ${cwd}/net/dl_file.sh ${URL}

# Checksum the URL
[[ ${SIG} != "" ]] && source ${cwd}/net/checksum.sh

# Extract the image file/archive
source ${cwd}/fs/extract_rootfs.sh

# Apply chroot configurations
source ${cwd}/fs/chroot.sh

# Pack the image
source ${cwd}/fs/makeimg.sh

# Convert to hekate format
[[ ${hekate} == "true" ]] && source ${cwd}/fs/hekate.sh
