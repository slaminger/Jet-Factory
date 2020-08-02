#!/usr/bin/bash
# ENTRYPOINT.SH : Manages options, and launches the sub scripts.

unset options i
while IFS= read -r -d $'\0' f; do
  options[i++]="$f"
done < <(find configs/ -maxdepth 1 -type f -print0 )

# Select DISTRO from avalaible one in configs
select opt in "${options[@]}"; do
	source ${opt}
	img="${URL##*/}"
break;
done

usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
	echo " --android            	Download setup files anyway"
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
	--hekate)  hekate=true; shift ;;
    ? | -h | --help) usage; exit 0 ;;
    -- ) shift; break ;;
    esac
done

# Build android using PabloZaiden's docker repo
if [ ${ROM_NAME} != "" ]; then
	mkdir -p ./android/lineage
	docker run --rm -ti -e ROM_NAME=${ROM_NAME} -v "$PWD"/android:/root/android pablozaiden/switchroot-android-build:latest
	exit 0
fi

# Download the URL
./net/dl_file.sh ${URL}

# Checksum the URL
[[ ${SIG} != "" ]] && ./net/checksum.sh

# Extract the image file/archive
./fs/extract_rootfs.sh

# Apply chroot configurations
arch-chroot ${rootdir} `${CHROOT_CMD}`

# Pack the image
./makeimg.sh

# Convert to hekate format
[[ ${hekate} == "true" ]] && ./fs/hekate.sh