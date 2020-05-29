#!/bin/bash

# TODO : Port this to Go
# echo -e "/dev/mmcblk0p1	/boot	vfat	rw,relatime	0	2\n" >> "${build_dir}/etc/fstab"
# sed -r -i 's/^HOOKS=((.*))$/HOOKS=(\1 resize-rootfs)/' "${build_dir}/etc/mkinitcpio.conf"
# 7z a "SWR-${img%.*}.7z" ${build_dir}/{bootloader,switchroot}

# Prepare root filesystem in Docker
docker run --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -ti -v /root/${DISTRO}:/root/${DISTRO} -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0 ./jetfactory -prepare -distro=${DISTRO}

docker run --rm --privileged multiarch/qemu-user-static:register --reset

# Chroot in filesystem and install packages
docker run --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -ti -v /dev:/dev -v /root/${DISTRO}:/root/${DISTRO} -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0 ./jetfactory -chroot -distro=${DISTRO}

# Make the final installable file
docker run --security-opt apparmor:unconfined --cap-add SYS_ADMIN --privileged --rm -ti -v /root/${DISTRO}:/root/${DISTRO} -v /var/run/docker.sock:/var/run/docker.sock alizkan/jet-factory:1.0.0 ./jetfactory -image -distro=${DISTRO}