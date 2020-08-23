#!/bin/bash
# HEKATE.SH Install hekate to a rootfs image and pack to hekate format
hekate_version=5.3.2
nyx_version=0.9.3
hekate_url=https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip
hekate_zip=${hekate_url##*/}
hekate_bin=hekate_ctcaer_${hekate_version}.bin

# Download hekate
wget -nc -q --show-progress ${hekate_url} -P "${out}/downloadedFiles/"

# Extract hekate
7z x "${out}/downloadedFiles/${hekate_zip}" -o"${out}/downloadedFiles/"

# Upload hekate payload using libguestfs
cp "${out}/downloadedFiles/${hekate_bin}" "/mnt/${NAME}_tmp_mnt/lib/firmware/reboot_payload.bin"

# Remove unneeded
rm "${out}/downloadedFiles/${hekate_zip}" "${out}/downloadedFiles/${hekate_bin}"

# Create switchroot install folder
mkdir -p "${out}/downloadedFiles/switchroot/install/"

# Unmount image
umount "/mnt/${NAME}_tmp_mnt"

# Split parts to output directory
split -b4290772992 --numeric-suffixes=0 "${guestfs_img}" "${out}/downloadedFiles/switchroot/install/l4t."

# 7zip the folder
7z a "SWR-${NAME}.7z" "${out}/downloadedFiles/bootloader"
7z a "SWR-${NAME}.7z" "${out}/downloadedFiles/switchroot"

# Clean hekate files and image
rm "${guestfs_img}"
rm -r "${out}/downloadedFiles/bootloader" "${out}/downloadedFiles/switchroot"
