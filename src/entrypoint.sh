#!/bin/bash
# ENTRYPOINT.SH : Manages options, and launches the sub scripts.

# If in docker, use the volume for linux
if [[ -f /.dockerenv ]]; then
	out="/root/linux"
# Else get the last argument passed to the script
else
	out=$(realpath ${@:$#})
fi

# Check if it's a valid path
if [[ ! -d ${out} ]]; then
	echo "Not a valid directory! Exiting.."
	exit 1
fi

# Store the script directory
cwd=$(dirname "$(readlink -f "$0")")

# Read configs dir for config files
unset options i
while IFS= read -r -d $'\0' f; do
  options[i++]="$f"
done < <(find $(dirname ${cwd})/configs/ -maxdepth 1 -type f -print0 )

if [[ ${DISTRO} == "" ]]; then
	echo "Select a configuration: "
	# Select DISTRO from configs files
	select opt in "${options[@]}"; do
		source ${opt}
		export $(cut -d= -f1 ${opt})
		img="${URL##*/}"
	break;
	done
else
	if [[ ${options[@]} =~ ${DISTRO} ]]; then
		source $(dirname ${cwd})/configs/${DISTRO}
		export $(cut -d= -f1 $(dirname ${cwd})/configs/${DISTRO})
		img="${URL##*/}"
	else
		echo "${DISTRO} couldn't be found in the configs/ directory! Exiting now..."
		exit 1
	fi
fi


echo "Preparing build directory..."
cd ${out}
mkdir -p ${out}/{${NAME},downloadedFiles}

echo "Adding executable bit to the scripts..."
chmod +x ${cwd}/{net,fs}/* $(dirname ${cwd})/configs/examples/*

echo "Downloading necessary files..."
source ${cwd}/net/dl_file.sh ${URL}

if [[ ${SIG} != "" ]]; then
	echo "Verifying file integrity..."
	source ${cwd}/net/checksum.sh
fi

echo "Extracting and preparing for chroot..."
source ${cwd}/fs/extract_rootfs.sh

echo "Chrooting..."
source ${cwd}/fs/chroot.sh

echo "Creating image file..."
source ${cwd}/fs/makeimg.sh

# Convert to hekate format
if [[ ${HEKATE} == "true" ]]; then
	echo "Creating hekate installable partition..."
	source ${cwd}/fs/hekate.sh
fi

echo "Done !"
