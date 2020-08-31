#!/bin/bash
# HEKATE.SH Install hekate to a rootfs image and pack to hekate format
hekate_version=5.3.2
nyx_version=0.9.3
hekate_url="https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip"
hekate_zip="${hekate_url##*/}"
hekate_bin="hekate_ctcaer_${hekate_version}.bin"
zip_final="switchroot-${DISTRO}.7z"
modules_dir="${out}/${NAME}/"
[[ -L "${modules_dir}/lib" ]] && modules_dir="${out}/${NAME}/usr/"

# Apply ext label
[[ -n "${HEKATE_ID}" ]] && echo -e "\nAssigning e2label: ${HEKATE_ID}\n" && \
	e2label "${guestfs_img}" "${HEKATE_ID}"

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

# Create 4MB aligned image
virt-make-fs --type=ext4 --format=raw --size=+$((512+${aligned_size_extra}))M "${out}/${NAME}/" ${guestfs_img}

# Split parts to output directory
split -b4290772992 --numeric-suffixes=0 "${guestfs_img}" "${out}/downloadedFiles/switchroot/install/l4t."

# 7zip the folder
7z a "${zip_final}" "${out}/downloadedFiles/switchroot/install/"

# Clean hekate files and image
rm -r "${out}/${guestfs_img}" "${out}/downloadedFiles/bootloader" "${out}/downloadedFiles/switchroot"
