#!/usr/bin/bash
# MAKEIMG.SH : Create rootfs image file

# Allocate size
size=$(du -hs -BM "${img}" | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M

# Create 4MB aligned image
dd if=/dev/zero of="${guestfs_img}" bs=1 count=0 seek=${size}

# Copy rootfs folder content using libguestfs to the image
# TODO