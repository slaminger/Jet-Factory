#!/bin/bash

# mv "${build_dir}/${hekate_bin}" "${build_dir}/lib/firmware/reboot_payload.bin"
# echo -e "/dev/mmcblk0p1	/boot	vfat	rw,relatime	0	2\n" >> "${build_dir}/etc/fstab"
# sed -r -i 's/^HOOKS=((.*))$/HOOKS=(\1 resize-rootfs)/' "${build_dir}/etc/mkinitcpio.conf"

findPkgManager() {
	# Source : https://ilhicas.com/2018/08/08/bash-script-to-install-packages-multiple-os.html
	declare -A osInfo;
	osInfo[/etc/debian_version]="apt-get update -y && apt-get install -y"
	osInfo[/etc/alpine-release]="apk"
	osInfo[/etc/centos-release]="yum update && yum install"
	osInfo[/etc/fedora-release]="dnf update && dnf install"
	osInfo[/etc/arch-release]="pacman -Syu"

	for f in ${!osInfo[@]}
	do
		if [[ -f $f ]];then
			package_manager=${osInfo[$f]}
		fi
	done

	echo $package_manager
}

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
docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0  ./jet-factory -prepare -distro=${DISTRO}
# docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jet-factory -configs -distro=${DISTRO}
# docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jet-factory -packages -distro=${DISTRO}
# docker run --privileged --rm -ti -v /var/run/docker.sock:/var/run/docker.sock azkali/jet-factory:1.0.0 ./jet-factory -image -distro=${DISTRO}