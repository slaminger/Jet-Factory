#!/usr/bin/bash
# MAKEIMG.SH : Create rootfs image file

# Allocate size
size=$(du -hs -BM "${rootdir}" | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M

# Create 4MB aligned image
dd if=/dev/zero of="${guestfs_img}" bs=1 count=0 seek=${size}

# EXtract tar content to image
virt-tar-in -a ${guestfs_img} ${tar_out} /