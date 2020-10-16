#!/bin/bash
set -e

# Variables :

# Check if output path is a valid path
out=$(realpath "${@:$#}")
# Store the script directory
cwd=$(dirname "$(readlink -f "$0")")
# Supported image file format
images_format=(".raw .img .iso")
# Output name of image file
guestfs_img="switchroot-${DISTRO}.img"

# Hekate Specific :

# Output name of hekate build 
zip_final="switchroot-${DISTRO}.7z"
hekate_version=5.3.3
nyx_version=0.9.4
hekate_url="https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip"
hekate_zip="${hekate_url##*/}"
hekate_bin="hekate_ctcaer_${hekate_version}.bin"

# Helper functions :

check() {
	if [[ ! -d "${out}" ]]; then
		echo "${out} is not a valid directory! Exiting.."
		exit 1
	fi

	if [[ ! -e "$(dirname "${cwd}")"/configs/"${DEVICE}" ]]; then
		echo "No directory for device: ${DEVICE} found in configs directory..."
		exit 1
	fi

	# Read configs dir for config files
	distro_avalaible=("$(dirname "${cwd}")"/configs/"${DEVICE}"/*)

	for distro_found in "${distro_avalaible[@]}"; do
		if [[ ${DISTRO} == "${distro_found##*/}" ]]; then
			set -a && . "${distro_found}" && set +a
			export img="${URL##*/}"
			break
		fi
	done

	if [[ -z "${img}" ]]; then
	    echo "${DISTRO} couldn't be found in the config directory! Exiting now..."
	    exit 1
	fi

}

hashsum() {
	if [[ -n "${SIG}" ]]; then
		echo -e "\nVerifying file integrity...\n" && \
		
		# Cut sig name from SIG URL
		img_sig="${SIG##*/}"

		# Download checksm if avalaible Check file integrity
		wget -q --show-progress "${SIG}" -P "${out}/downloadedFiles/"

		# Checksum
		if [[ "${SIG}" =~ .md5 ]]; then
			md5sum --status -c "${out}/downloadedFiles/${img_sig}"
		else
			shasum --status -c "${out}/downloadedFiles/${img_sig}"
		fi
	fi
}

# Core functions :

prepare() {
	if [[ -d "${out}/${NAME}" ]]; then
		echo -e "Cleaning previous build directory...\n"
		rm -rf "${out}/${NAME}"
	fi

	echo -e "Preparing build directory...\n"
	cd "${out}"
	mkdir -p "${out}/${NAME}" "${out}/downloadedFiles"

	echo -e "Adding executable bit to the scripts...\n"
	chmod +x "${cwd}"/{net,fs}/* "$(dirname "${cwd}")"/configs/"${DEVICE}"/files/*

	echo -e "Downloading necessary files...\n"
	if [[ ! -e "${out}/downloadedFiles/${img%.*}" ]]; then
		if [[ -n "${URL}" ]]; then
			wget -q -nc --show-progress "${URL}" -P "${out}/downloadedFiles/"
		else
			echo "No URL found !";exit 1;
		fi
	fi
}

extrct_rootfs() {
	echo -e "\nExtracting and preparing for chroot...\n"

	img="${out}/downloadedFiles/${img}"
	for format in $images_format; do
		[[ ${img%.*} = *$format ]] && is_iso=1 && break;
	done

	# Handle rootfs extraction
	if [[ "${img}" = *.tbz2 ]]; then
		tar xpf --xattrs-include='*.*' --numeric-owner "${img}" -C "${out}/${NAME}"
	elif [[ "${img}" =~ .tar ]]; then
		bsdtar xpf --xattrs-include='*.*' --numeric-owner "${img}" -C "${out}/${NAME}"
	elif [[ -n ${is_iso} ]]; then
		# Handle xz compressed images
		if [[ "$(file -b --mime-type "${img}")" == "application/x-xz" ]]; then
			[[ ! -e "${img%.*}" ]] && unxz "${img}"
			img="${img%.*}"
		fi

		echo -e "Scanning image file for rootfs partition\n"
		rootfs="$(guestfish -a "${img}" launch : inspect-os)"

		echo -e "Extracting partition from image file. This will take a while...\n"
		virt-copy-out -a "${img}" -m "${rootfs}" / "${out}/${NAME}"
	else
		echo -e "Unrecognized format.\n"
		exit 1
	fi
}

chroot() {
	echo -e "\nChrooting...\n"
	AARCH=$(guestfish -a ${img} launch : inspect-os : inspect-get-arch | tail -1)
	[[ "$(uname -m)" == "${AARCH}" ]] && same_arch=1

	# Check if architecture is already registered in the .lock file
	if [[ -f "${out}/.lock" ]]; then
		lock="$(grep -xF "${AARCH}" "${out}/.lock")"
	else
		echo "${AARCH} 1" > "${out}/.lock"
	fi

	# If an architecture is  already registered increment the counter
	if [[ -n ${lock} ]]; then
		# Get current lock count
		lock_count=$(sed 's/'${AARCH}' //g' "${out}/.lock")

		# Increment lock count in lock file
		sed -i 's/'${AARCH}' '${lock_count}'/'${AARCH}' '$((lock_count+1))'/g' "${out}/.lock"

		# Increment lock count variable
		lock_count=$((lock_count+1))
	fi

	if [[ -z ${same_arch} && -z ${lock} ]]; then
		if [ ! -f /proc/sys/fs/binfmt_misc/register ]; then
			if ! mount binfmt_misc -t binfmt_misc /proc/sys/fs/binfmt_misc; then
			exit 1
		    fi
		fi
		
		if [[ ! -e "/proc/sys/fs/binfmt_misc/qemu-${AARCH}" ]]; then
			wget -L -q -nc --show-progress https://raw.githubusercontent.com/dbhi/qus/main/register.sh -P "${out}/downloadedFiles/"
			chmod +x "${out}/downloadedFiles/register.sh"
			"${out}/downloadedFiles/register.sh" -s -- -p "${AARCH}"
			cp "/usr/bin/qemu-${AARCH}-static" "${out}/${NAME}/usr/bin/"
		fi
	fi

	# Mount bind chroot dir
	mount --bind "${out}/${NAME}" "${out}/${NAME}"

	# Mounts switchroot folder as boot folder if a hekate ID is given
	if [[ -n ${HEKATE_ID} ]]; then
		if [[ -e "${out}/switchroot/${DISTRO}" ]]; then
			mount --bind "${out}/switchroot/${DISTRO}" "${out}/${NAME}/boot/"
		fi

		if [ -e "${out}/switchroot/${DISTRO}/update.tar.gz" ]; then
			tar xhpf "${out}/switchroot/${DISTRO}/update.tar.gz" -C "${out}/${NAME}"
		fi

	fi

	# Add cache dir configuration
	if [[ -n "$CACHE_DIR" ]]; then
		mkdir "${out}/cache" &> /dev/null || true
		mount --bind "${out}/cache" "${out}/${NAME}/${CACHE_DIR}" || exit
	fi

	# Copy build script
	cp "$(dirname "${cwd}")/configs/${DEVICE}/files/${CHROOT_SCRIPT}" "${out}/${NAME}"

	# Handle resolv.conf
	cp --remove-destination --dereference /etc/resolv.conf "${out}/${NAME}/etc/resolv.conf"

	# Actual chroot
	arch-chroot "${out}/${NAME}" /bin/bash /"${CHROOT_SCRIPT}"

	# Unmount switchroot boot dir
	[[ -n ${HEKATE_ID} ]] && umount -l "${out}/${NAME}/boot"

	# Unmount chroot dir
	umount -l "${out}/${NAME}"

	# Check lock status
	if [[ -z ${lock} ]]; then
		# Get current lock count
		lock_count=$(sed 's/'${AARCH}' //g' "${out}/.lock")

		# If the current instance is the only one left for this binary, remove it
		if [[ ${lock_count} = 1 ]]; then
			# Remove lock on architecture
			sed -i '/'${AARCH}'*/d' "${out}/.lock"

			# Unregister binary if it wasn't set on script launch
			"${out}/downloadedFiles/register.sh" -- -r -p ${AARCH}
		else
			# Decrement lock count in lock file
			sed -i 's/'${AARCH}' '${lock_count}'/'${AARCH}' '$((lock_count-1))'/g' "${out}/.lock"
		fi
	fi

	# Remove lock file if empty, meaning no more instance is running.
	[[ ! -s "${out}/.lock" ]] && rm -rf "${out}/.lock"

	# Clean qemu emulation files
	if [[ -z ${same_arch} ]]; then
		rm -rf "${out}/${NAME}/usr/bin/qemu-${AARCH}-static" "${out}/downloadedFiles/register.sh"
	fi

	[[ -n "$CACHE_DIR" ]] && umount -l "${out}/cache"

	rm -rf "${out}/${NAME}/${CHROOT_SCRIPT}" 
}

create_image() {
	echo -e "\nCreating image file...\n"

	# Clean previously made image file or 7zip
	[[ -f "${guestfs_img}" ]] && rm -rf "${guestfs_img}"
	[[ -f "${zip_final}" ]] && rm -rf "${zip_final}"

	if [[ -n "${HEKATE_ID}" ]]; then
		modules_dir="${out}/${NAME}/"
		[[ -L "${modules_dir}/lib" && -d "${modules_dir}/lib" ]] && modules_dir="${out}/${NAME}/usr/"

		# Download hekate
		wget -nc -q --show-progress ${hekate_url} -P "${out}/downloadedFiles/"

		# Extract hekate
		7z x "${out}/downloadedFiles/${hekate_zip}" ${hekate_bin}

		# Copy hekate bin
		mv "${hekate_bin}" "${modules_dir}/lib/firmware/"

		# Remove unneeded
		rm "${out}/downloadedFiles/${hekate_zip}" 
	fi

	# Create image
	virt-make-fs --type=ext4 --format=raw --size=+512MB "${out}/${NAME}/" ${guestfs_img}

	# Zerofree the image produced
	zerofree -n ${guestfs_img}

	# Apply ext label
	if [[ -n "${HEKATE_ID}" ]]; then
		echo -e "\nAssigning e2label: ${HEKATE_ID}\n"
		e2label "${guestfs_img}" "${HEKATE_ID}"
	fi

	# Convert to hekate format or create image
	if [[ "${HEKATE}" == "true" ]]; then
		echo "Creating hekate installable partition..."

		# Create switchroot install folder
		mkdir -p "${out}/downloadedFiles/switchroot/install/"

		# Get build directory size
		size="$(du -b -s "${guestfs_img}" | awk '{print int($1);}')"

		# Alignement adjust to 4MB
		aligned_size=$(((${size} + (4194304-1)) & ~(4194304-1)))

		# Check if image needs alignement
		align_check=$((${aligned_size} - ${size}))

		# Align part if necessary
		[[ ${align_check} -ne 0 ]] && dd if=/dev/zero bs=1 count=${align_check} >> ${guestfs_img}

		# Split parts to output directory
		split -b4290772992 --numeric-suffixes=0 "${guestfs_img}" "${out}/downloadedFiles/switchroot/install/l4t."

		# 7zip the folder
		7z a "${zip_final}" "${out}/downloadedFiles/switchroot/"

		# Clean hekate files and image
		rm -rf "${out}/${NAME}/" "${out}/${guestfs_img}" "${out}/downloadedFiles/bootloader" "${out}/downloadedFiles/switchroot"

		echo -e "\nDone ! Hekate flashable 7zip resides in ${out}/${zip_final}"
	else
		# Clean unneeded files
		rm -rf "${out}/${NAME}/"
		echo -e "\nDone ! Image resides in ${out}/${guestfs_img}"
	fi
}
