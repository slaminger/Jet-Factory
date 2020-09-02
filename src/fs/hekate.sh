#!/bin/bash
# HEKATE.SH Install hekate to a rootfs image and pack to hekate format
hekate_version=5.3.2
nyx_version=0.9.3
hekate_url="https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip"
hekate_zip="${hekate_url##*/}"
hekate_bin="hekate_ctcaer_${hekate_version}.bin"
zip_final="switchroot-${DISTRO}.7z"
modules_dir="${out}/${NAME}/"
[[ -L "${modules_dir}/lib" && -d "${modules_dir}/lib" ]] && modules_dir="${out}/${NAME}/usr/"

[[ -f "${zip_final}" ]] && rm -rf "${zip_final}"

# Download hekate
wget -nc -q --show-progress ${hekate_url} -P "${out}/downloadedFiles/"

# Extract hekate
7z x "${out}/downloadedFiles/${hekate_zip}" -o"${out}/downloadedFiles/"

# Upload hekate payload
cp "${out}/downloadedFiles/${hekate_bin}" "${modules_dir}/lib/firmware/reboot_payload.bin"

# Remove unneeded
rm "${out}/downloadedFiles/${hekate_zip}" "${out}/downloadedFiles/${hekate_bin}"

# Create switchroot install folder
mkdir -p "${out}/downloadedFiles/switchroot/install/"

# Create image
virt-make-fs --type=ext4 --format=raw --size=+512MB "${out}/${NAME}/" ${guestfs_img}

zerofree -n ${guestfs_img}

# Get build directory size
size="$(du -b -s "${guestfs_img}" | awk '{print int($1);}')"

# Alignement adjust to 4MB
aligned_size=$(((${size} + (4194304-1)) & ~(4194304-1)))

# Check if image needs alignement
align_check=$((${aligned_size} - ${size}))

# Apply ext label
if [[ -n "${HEKATE_ID}" ]]; then
	echo -e "\nAssigning e2label: ${HEKATE_ID}\n"
	e2label "${guestfs_img}" "${HEKATE_ID}"
fi

# Align part if necessary
[[ ${align_check} -ne 0 ]] && dd if=/dev/zero bs=1 count=${align_check} >> ${guestfs_img}

# Split parts to output directory
split -b4290772992 --numeric-suffixes=0 "${guestfs_img}" "${out}/downloadedFiles/switchroot/install/l4t."

# 7zip the folder
7z a "${zip_final}" "${out}/downloadedFiles/switchroot/"

# Clean hekate files and image
rm -rf "${out}/${guestfs_img}" "${out}/downloadedFiles/bootloader" "${out}/downloadedFiles/switchroot"
