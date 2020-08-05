#!/usr/bin/bash
# MAKEIMG.SH : Create rootfs image file
tar_tmp=${NAME}.tar

# Set Image name
guestfs_img="SWR-${NAME}.img"

# Allocate size
size=$(du -hs -BM "${out}/${NAME}/" | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M

# Create 4MB aligned image
dd if=/dev/zero of=${guestfs_img} bs=1 count=0 seek=${size}

# Create ext4 partition
mkfs.ext4 ${guestfs_img}

# Create temporary tar archive
tar cf ${tar_tmp} "${out}/${NAME}/"

# EXtract tar content to image
tar xOf  ${tar_tmp} | dd of=${guestfs_img} bs=1M

# Remove unneeded files
rm -r ${tar_tmp} "${out}/${NAME}/"