#!/bin/bash
build_dir=$1

echo "Post chroot cleaning..."
umount "${build_dir}/boot/" && umount ${build_dir}
rm ${build_dir}/usr/bin/qemu-aarch64-static