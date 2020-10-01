#!/bin/bash
# ENTRYPOINT.SH : Manages distro_avalaible, and launches the sub scripts.
set -e

# Check if output path is a valid path
out=$(realpath "${@:$#}")

[[ ! -d "${out}" ]] && \
	echo "${out} is not a valid directory! Exiting.." && exit 1

# Store the script directory
cwd=$(dirname "$(readlink -f "$0")")

[[ ! -e "$(dirname "${cwd}")"/configs/"${DEVICE}" ]] && \
	echo "No directory for device: ${DEVICE} found in configs directory..." && exit 1

# Read configs dir for config files
distro_avalaible=("$(dirname "${cwd}")"/configs/"${DEVICE}"/*)

for distro_found in "${distro_avalaible[@]}"; do
	if [[ ${DISTRO} == "${distro_found##*/}" ]]; then
		set -a && . "${distro_found}" && set +a
		export img="${URL##*/}"
		break
	fi
done

[[ -z "${img}" ]] && \
    echo "${DISTRO} couldn't be found in the config directory! Exiting now..." && exit 1

if [[ -d "${out}/${NAME}" ]]; then
	echo -e "\nCleaning previous build directory...\n"
	rm -rf "${out}/${NAME}"
fi

echo -e "\nPreparing build directory...\n"
cd "${out}" || exit
mkdir -p "${out}"/{"${NAME}",downloadedFiles}

echo -e "\nAdding executable bit to the scripts...\n"
chmod +x "${cwd}"/{net,fs}/* "$(dirname "${cwd}")"/configs/"${DEVICE}"/files/*

echo -e "\nDownloading necessary files...\n"
source "${cwd}/net/dl_file.sh" "${URL}"

[[ -n "${SIG}" ]] && echo -e "\nVerifying file integrity...\n" && \
	source "${cwd}/net/checksum.sh"

echo -e "\nExtracting and preparing for chroot...\n"
source "${cwd}/fs/extract_rootfs.sh"

echo -e "\nChrooting...\n"
source "${cwd}/fs/chroot.sh"

echo -e "\nCreating image file...\n"
source "${cwd}/fs/makeimg.sh"
echo -e "\nDone !\n"
