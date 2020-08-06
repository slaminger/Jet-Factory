#!/usr/bin/bash
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
mkdir -p /tmp/${NAME}_tmp_mnt

# Mount the disk image
mount ${guestfs_img} /tmp/${NAME}_tmp_mnt

# Copy files
cp -a ${out}/${NAME}/* /tmp/${NAME}_tmp_mnt

# Unmount image
umount /tmp/${NAME}_tmp_mnt

# Remove unneeded files
rm -r "${out}/${NAME}/"