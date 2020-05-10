#!/bin/bash
build_dir=$1

echo "Pre chroot setup..."
cp /usr/bin/qemu-aarch64-static ${build_dir}/usr/bin/
mount --bind ${build_dir} ${build_dir} &&
mount --bind  "${build_dir}/bootloader/" "${build_dir}/boot/"