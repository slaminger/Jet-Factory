#!/bin/bash
# MAKEIMG.SH : Create rootfs image file
tar_tmp=${NAME}.tar

# Set Image name
guestfs_img="SWR-${NAME}.img"

# Clean previously made image file
[[ -f ${guestfs_img} ]] && rm ${guestfs_img}

# Allocate size
size=$(du -hs -BM "${out}/${NAME}/" | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M

# Create 4MB aligned image
dd if=/dev/zero of=${guestfs_img} bs=1 count=0 seek=${size}

# Create ext4 partition
mkfs.ext4 -F ${guestfs_img}

# Create tmp directroy
mkdir -p "/mnt/${NAME}_tmp_mnt"

# Mount the disk image
mount ${guestfs_img} "/mnt/${NAME}_tmp_mnt"

# Copy files
cp -a ${out}/${NAME}/* "/mnt/${NAME}_tmp_mnt"

# Convert to hekate format
if [[ ${HEKATE} == "true" ]]; then
	echo "Creating hekate installable partition..."
	source ${cwd}/fs/hekate.sh
else
	# Or unmount image
	umount "/mnt/${NAME}_tmp_mnt"
fi

# Remove unneeded files
rm -r "${out}/${NAME}/" "/mnt/${NAME}_tmp_mnt"
