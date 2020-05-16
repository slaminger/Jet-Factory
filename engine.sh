#!/bin/bash
# echo -e "/dev/mmcblk0p1	/boot	vfat	rw,relatime	0	2\n" >> "${build_dir}/etc/fstab"

# sed -r -i 's/^HOOKS=((.*))$/HOOKS=(\1 resize-rootfs)/' "${build_dir}/etc/mkinitcpio.conf"

CreateImage() {
	echo "Compressing hekate folder..."
	if [[ "${dl_dir}/${img%.*}" =~ '.tar' ]]; then 
		7z a "SWR-${img%%.*}.7z" ${build_dir}/{bootloader,switchroot}
	else
		7z a "SWR-${img%.*}.7z" ${build_dir}/{bootloader,switchroot}
	fi
}

# Actual script
docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jetfactory -prepare -distro=${DISTRO}
# docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jetfactory -configs -distro=${DISTRO}
# docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jetfactory -packages -distro=${DISTRO}
# docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jetfactory -image -distro=${DISTRO}