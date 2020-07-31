#!/usr/bin/bash
# HEKATE.SH Install hekate to a rootfs image and pack to hekate format
hekate_version=5.3.2
nyx_version=0.9.3
hekate_url=https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip
hekate_zip=${hekate_url##*/}
hekate_bin=hekate_ctcaer_${hekate_version}.bin

# Download hekate
wget -P "$1" -q --show-progress ${hekate_url} -O "$1/${hekate_zip}"

# Extract hekate
unzip -q -o "$1/${hekate_zip}" -d "$1"

# Upload hekate payload using libguestfs
virt-copy-in -a ${img} "$1/${hekate_bin}" /usr/lib/firmware/reboot_payload.bin

# Remove unneeded zip
rm "$1/${hekate_zip}"

# Split parts to output directory
split -b4290772992 --numeric-suffixes=0 "${img}" "$1/switchroot/install/l4t."

# 7zip the folder
7z a "SWR-${img}.7z" $1/*