#!/bin/env bash

# Setup variables
docker=true
staging=false
hekate=false

# Folders
cwd="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
build_dir="${cwd}/build"
dl_dir="${cwd}/dl"

# Android version specific variables
selection="$(echo $1 | tr '[:upper:]' '[:lower:]')"
img_url="$(head -1 ${cwd}/${selection}/urls)"
img_sig_url="${img_url}.md5"
img="SWR-"${img_url##*/}""
img_sig="${img_sig_url##*/}"
validate_command="md5sum --status -c "${img_sig}""

# Hekate files
hekate_version=5.2.0
nyx_version=0.9.0
hekate_url=https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip
hekate_zip=${hekate_url##*/}
hekate_bin=hekate_ctcaer_${hekate_version}.bin

SetVersion() {
	if [[ ${selection} == "pie" ]]; then
		img_sig_url=
	elif [[ ${selection} == "oreo" ]]; then
		img_sig_url=
	else
		echo "$0: invalid distro option: $1"
		usage
		exit 1
	fi
}

GetImgFiles() {
	# cd into download directory
	cd ${dl_dir}

	# Download file if it doesn't exist, or is forced to download.
	if [[ ! -f ${img} || $1 == "force" ]]; then 
		wget -q --show-progress ${img_url} -O ${img}
	else
		echo "Image exists!"
	fi
	
	# Download signature file
	echo "Downloading signature file..."
	wget -q --show-progress ${img_sig_url} -O ${img_sig}
	
	# Check image against signature
	echo "Validating image..."
	$validate_command
	if [[ $? != "0" ]]; then
		echo "Image doesn't match signature, re-downloading..."
		GetImgFiles force
	else
		echo "Signature check passed!"
	fi
}

Main() {
	# TODO : Build android
	git clone https://github.com/PabloZaiden/switchroot-android-
	cd switchroot-android-
	./build.sh "${@: -1}"
}

# Parse arguments
options=$(getopt -n $0 -o dfhns --long force,hekate,no-docker,staging,distro:,help -- "$@")

# Check for errors in arguments or if no name was provided
if [[ $? != "0" ]] || [[ "${@: -1}" =~ options ]]; then usage; exit 1; fi

# Evaluate arguments
eval set -- "$options"
while true; do
    case "$1" in
	-d | --docker) docker=true; shift ;;
    -s | --staging) staging=true; shift ;;
	--hekate) hekate=true; shift ;;
    ? | -h | --help) usage; exit 0 ;;
    -- ) shift; break ;;
    esac
done

Main