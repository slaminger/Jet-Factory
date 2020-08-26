#!/bin/bash
# ENTRYPOINT.SH : Manages distro_avalaible, and launches the sub scripts.

# Check if output path is a valid path
out=$(realpath "${@:$#}")

[[ ! -d "${out}" ]] && \
	echo "${out} is not a valid directory! Exiting.." && exit 1

# Store the script directory
cwd=$(dirname "$(readlink -f "$0")")

[[ ! -e "${cwd}"/configs/"${DEVICE}" ]] && \
	echo "No directory for device: ${DEVICE} found in configs directory..." && exit 1

# Read configs dir for config files
unset distro_avalaible i
while IFS= read -r -d $'\0' f; do
  distro_avalaible[i++]="$f"
done < <(find "$(dirname "${cwd}")"/configs/"${DEVICE}" -maxdepth 1 -type f -print0)

for distro_found in "${distro_avalaible[@]}"
do
	if [[ "${DISTRO}" =~ ${distro_found}  ]]; then
		# shellcheck disable=SC1090 disable=SC1091
		source "$(dirname "${cwd}")/configs/${DEVICE}/${distro_found}"
		export "$(cut -d= -f1 "$(dirname "${cwd}")"/configs/"${DEVICE}"/"${distro_found}")"
		export img="${URL##*/}"
	else
		echo "${DISTRO} couldn't be found in the config directory! Exiting now..."
		exit 1
	fi
done

echo "Preparing build directory..."
cd "${out}" || exit
mkdir -p "${out}"/{"${NAME}",downloadedFiles}

echo "Adding executable bit to the scripts..."
chmod +x "${cwd}"/{net,fs}/* "$(dirname "${cwd}")"/configs/examples/*

echo "Downloading necessary files..."
# shellcheck source=src/net/dl_file.sh disable=SC1091
source "${cwd}/net/dl_file.sh" "${URL}"

# shellcheck source=src/net/checksum.sh disable=SC1091
[[ -n "${SIG}" ]] && echo "Verifying file integrity..." && \
	source "${cwd}/net/checksum.sh"

echo "Extracting and preparing for chroot..."
# shellcheck source=src/fs/extract_rootfs.sh disable=SC1091
source "${cwd}/fs/extract_rootfs.sh"

echo "Chrooting..."
# shellcheck source=src/fs/chroot.sh disable=SC1091
source "${cwd}/fs/chroot.sh"

echo "Creating image file..."
# shellcheck source=src/fs/makeimg.sh disable=SC1091
source "${cwd}/fs/makeimg.sh"
echo "Done !"
