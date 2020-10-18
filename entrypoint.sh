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
guestfs_img="${out}/switchroot-${DISTRO}.img"

# Hekate Specific :

# Output name of target hekate build 
zip_final="${out}/switchroot-${DISTRO}.7z"
hekate_version=5.3.3
nyx_version=0.9.4
hekate_url="https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip"
hekate_zip="${hekate_url##*/}"
hekate_bin="hekate_ctcaer_${hekate_version}.bin"

# Helper functions :

get_file() {
	img="${URL##*/}"

	echo -e "Downloading necessary files...\n"
	if [[ ! -e "${out}/downloadedFiles/${img%.*}" ]]; then
		wget -q -nc --show-progress "${URL}" -P "${out}/downloadedFiles/"
	fi
}

hashsum() {
	if [[ -n "${SIG}" ]]; then
		echo -e "Verifying file integrity...\n" && \
		
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

cleanup() {
	echo -e "Cleaning build files...\n"

	# Unmount cache when done
	if [[ "$(mountpoint -q ${out}/cache)" ]]; then
		umount "${out}/cache"
	fi

	# Unmount chroot dir
	if [[ -d "${build_dir}" ]]; then
		if [[ "$(mountpoint -q ${build_dir})" ]]; then
			umount -R "${build_dir}"
		fi

		rm -rf "${build_dir}"
	fi
}

# Core functions :

prepare() {
	echo -e "Checking variables\n"

	if [[ ! -d "${out}" ]]; then
		echo "${out} is not a valid directory! Exiting.."
		exit 1
	fi

	if [[ -z "${DEVICE}" ]]; then
		echo "No device specified. Exiting !"
		exit 1
	fi

	if [[ ! -e "${cwd}/configs/${DEVICE}" ]]; then
		echo "No device name : ${DEVICE} could be found in config. Exiting !"
		exit 1
	fi

	# Read configs dir for config files
	distro_avalaible=(${cwd}/configs/${DEVICE}/*)

	for distro_found in "${distro_avalaible[@]}"; do
		if [[ ${DISTRO} == "${distro_found##*/}" ]]; then
			set -a && . "${distro_found}" && set +a
			break
		fi
	done

	if [[ -z "${URL}" ]]; then
		echo "No URL found. Exiting."
		exit 1
	fi

	if [[ -z "${CHROOT_SCRIPT}" ]]; then
	    echo "No CHROOT_SCRIPT found. Exiting !"
	    exit 1
	fi

	echo -e "Preparing build directory...\n"
	build_dir="${out}/${DEVICE}-${DISTRO}"
	mkdir -p "${build_dir}" "${out}/downloadedFiles"

	echo -e "Adding executable bit to the scripts...\n"
	chmod +x ${cwd}/configs/${DEVICE}/files/*.sh ${cwd}/*.sh
}

extract_rootfs() {
	echo -e "Extracting and preparing for chroot...\n"

	img="${out}/downloadedFiles/${img}"
	for format in $images_format; do
		if [[ ${img} = *$format ]]; then
			is_iso=1
			break
		elif [[ ${img%.*} = *$format ]]; then
			is_iso=1; img=${img%.*};
			break
		fi
	done

	if [[ "${img}" = *.tbz2 ]]; then
		tar xpf --xattrs-include='*.*' --numeric-owner "${img}" -C "${build_dir}"
	fi

	if [[ "${img}" =~ .tar ]]; then
		bsdtar xpf --xattrs-include='*.*' --numeric-owner "${img}" -C "${build_dir}"
	fi

	if [[ -n ${is_iso} ]]; then
		# Handle xz compressed images
		if [[ "$(file -b --mime-type "${img}")" == "application/x-xz" ]]; then
			[[ ! -e "${img%.*}" ]] && unxz "${img}"
		fi

		echo -e "Extracting partition from image file. This will take a while...\n"
		virt-copy-out -a "${img}" / "${build_dir}"
	fi
}

chroot_wrapper() {
	echo -e "Chrooting...\n"

	# Get root partition using guestfish
	rootfs=$(guestfish -a ${img} launch : inspect-os)

	# Get OS architecture using guestfish
	AARCH=$(guestfish -a ${img} launch : inspect-os : inspect-get-arch ${rootfs} | tail -1)

	# Check if target and host architecture are the same
	[[ "$(uname -m)" == "${AARCH}" ]] && same_arch=1

	# Register binary if the host and target CPU architectures differ
	# and no other instance of same CPU emulation is happening
	if [[ -z ${same_arch} ]]; then
		if [ ! -f /proc/sys/fs/binfmt_misc/register ]; then
			if ! mount binfmt_misc -t binfmt_misc /proc/sys/fs/binfmt_misc; then
			exit 1
		    fi
		fi
		
		if [[ ! -e "/proc/sys/fs/binfmt_misc/qemu-${AARCH}" ]]; then
			"${cwd}/register.sh" -s -- -p "${AARCH}"
		fi

		cp "/usr/bin/qemu-${AARCH}-static" "${build_dir}/usr/bin/"
	fi

	# Mount bind chroot dir
	mount --bind "${build_dir}" "${build_dir}"

	# Add cache dir configuration
	if [[ -n "${CACHE_DIR}" ]]; then
		mkdir -p "${out}/cache"
		mount --bind "${out}/cache" "${build_dir}/${CACHE_DIR}"
	fi

	# Mounts switchroot folder as boot folder if a hekate ID is given
	if [[ -n ${HEKATE_ID} ]]; then
		if [[ -e "${out}/switchroot/${DISTRO}" ]]; then
			mount --bind "${out}/switchroot/${DISTRO}" "${build_dir}/boot/"
		fi

		if [ -e "${out}/switchroot/${DISTRO}/update.tar.gz" ]; then
			tar xhpf "${out}/switchroot/${DISTRO}/update.tar.gz" -C "${build_dir}"
		fi
		
		if [ -e "${out}/switchroot/${DISTRO}/modules.tar.gz" ]; then
			tar xhpf "${out}/switchroot/${DISTRO}/modules.tar.gz" -C "${build_dir}"
		fi
	fi

	# Copy build script
	cp "${cwd}/configs/${DEVICE}/files/${CHROOT_SCRIPT}" "${build_dir}"

	# Handle resolv.conf
	cp --remove-destination /etc/resolv.conf "${build_dir}/etc/resolv.conf"

	# Create a lock with file descriptor: 200
	(
		flock -x 200
		
		# Ensure lock file is removed after the process ends
		trap "rm -rf ${out}/.lock-${AARCH}" 0

		# Actual chroot
		arch-chroot "${build_dir}" /bin/bash "/${CHROOT_SCRIPT}"
	) 200> "${out}/.lock-${AARCH}"

	if [[ -z $same_arch ]]; then
		rm -rf "${build_dir}/usr/bin/qemu-${AARCH}-static" \
			"${build_dir}/${CHROOT_SCRIPT}"
	fi
}

create_target() {
	echo -e "Creating image file. This will take a while...\n"

	if [[ -n "${HEKATE_ID}" ]]; then
		# Default payload storage would be /lib/
		modules_dir="${build_dir}"

		# If /lib is a symlink then it should be placed in /usr/lib/
		[[ -L "${modules_dir}/lib" && -d "${modules_dir}/lib" ]] && modules_dir="${build_dir}/usr/"

		# Download hekate release
		wget -nc -q --show-progress ${hekate_url} -P "${out}/downloadedFiles/"

		# Extract hekate bin from releas
		7z x "${out}/downloadedFiles/${hekate_zip}" ${hekate_bin}

		# Copy hekate bin to filesystem
		mv "${hekate_bin}" "${modules_dir}/lib/firmware/"

		# Remove unneeded
		rm "${out}/downloadedFiles/${hekate_zip}" 
	fi

	# Create image
	virt-make-fs --type=ext4 --format=raw --size=+512MB "${build_dir}" ${guestfs_img}

	# Zerofree the image produced
	zerofree -n ${guestfs_img}

	# Apply ext label
	if [[ -n "${HEKATE_ID}" ]]; then
		echo -e "Assigning e2label: ${HEKATE_ID}\n"
		e2label "${guestfs_img}" "${HEKATE_ID}"
	fi

	if [[ "${HEKATE}" = "true" ]]; then
		create_hekate_zip
		echo -e "Done ! Hekate flashable 7zip resides in ${out}/${zip_final}"
	else
		echo -e "Done ! Image resides in ${out}/${guestfs_img}"
	fi
}

create_hekate_zip() {
	# Convert to hekate format or create image
	echo -e "Creating hekate installable partition...\n"

	# Create switchroot install folder
	switchroot_dir="${build_dir}/switchroot"
	mkdir -p "${switchroot_dir}/install/"

	# Get build directory size
	size="$(du -b -s "${guestfs_img}" | awk '{print int($1);}')"

	# Alignement adjust to 4MB
	aligned_size=$(((${size} + (4194304-1)) & ~(4194304-1)))

	# Check if image needs alignement
	align_check=$((${aligned_size} - ${size}))

	# Align part if necessary
	if [[ ${align_check} -ne 0 ]]; then
		dd if=/dev/zero bs=1 count=${align_check} >> ${guestfs_img}
	fi

	# Split parts to output directory
	split -b4290772992 --numeric-suffixes=0 "${guestfs_img}" "${switchroot_dir}/install/l4t."

	# 7zip the folder
	7z a "${zip_final}" "${switchroot_dir}"

	# Clean image
	rm -rf "${guestfs_img}" 
}

# Cleanup on CTRL_C signal
trap cleanup INT

# Cleanup on EXIT signal
trap cleanup 0

cleanup
prepare
get_file
hashsum
extract_rootfs
chroot_wrapper
create_target
