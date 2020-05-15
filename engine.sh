#!/bin/bash
# mv "${build_dir}/${hekate_bin}" "${build_dir}/lib/firmware/reboot_payload.bin"
# echo -e "/dev/mmcblk0p1	/boot	vfat	rw,relatime	0	2\n" >> "${build_dir}/etc/fstab"

# sed -r -i 's/^HOOKS=((.*))$/HOOKS=(\1 resize-rootfs)/' "${build_dir}/etc/mkinitcpio.conf"

CreateImage() {
	size=$(du -hs -BM ${build_dir} | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M

	echo "Creating Hekate installable files..."
	split -b4290772992 --numeric-suffixes=0 ${img} l4t.
	
	echo "Compressing hekate folder..."
	if [[ "${dl_dir}/${img%.*}" =~ '.tar' ]]; then 
		7z a "SWR-${img%%.*}.7z" ${build_dir}/{bootloader,switchroot}
	else
		7z a "SWR-${img%.*}.7z" ${build_dir}/{bootloader,switchroot}
	fi
}

# Actual script
docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 sh -c "./jetfactory -prepare -distro=${DISTRO} && ls */*/*"
# docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jetfactory -configs -distro=${DISTRO}
# docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jetfactory -packages -distro=${DISTRO}
# docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jetfactory -image -distro=${DISTRO}