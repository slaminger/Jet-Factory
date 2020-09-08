#!/bin/bash
# MAKEIMG.SH : Create rootfs image file
zip_final="switchroot-${DISTRO}.7z"
guestfs_img="switchroot-${DISTRO}.img"

# Clean previously made image file or 7zip
[[ -f "${guestfs_img}" ]] && rm -rf "${guestfs_img}"
[[ -f "${zip_final}" ]] && rm -rf "${zip_final}"

if [[ -n "${HEKATE_ID}" ]]; then
	hekate_version=5.3.2
	nyx_version=0.9.3
	hekate_url="https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip"
	hekate_zip="${hekate_url##*/}"
	hekate_bin="hekate_ctcaer_${hekate_version}.bin"
	modules_dir="${out}/${NAME}/"
	[[ -L "${modules_dir}/lib" && -d "${modules_dir}/lib" ]] && modules_dir="${out}/${NAME}/usr/"

	# Download hekate
	wget -nc -q --show-progress ${hekate_url} -P "${out}/downloadedFiles/"

	# Extract hekate
	7z x "${out}/downloadedFiles/${hekate_zip}" ${hekate_bin} -o"${modules_dir}/lib/firmware/reboot_payload.bin"

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
	rm -rf "${out}/${guestfs_img}" "${out}/downloadedFiles/bootloader" "${out}/downloadedFiles/switchroot"
fi

# Clean unneeded files
rm -rf "${out}/${NAME}/"
