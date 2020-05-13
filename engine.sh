#!/bin/bash
### TODO : Port this to Go
PrepareFiles() {
	unzip ${hekate_zip} -d ${build_dir}

	if [[ ! -f ${dl_dir}/${img%.*} ]]; then
		cd ${build_dir}
		case ${img} in
			*.tar)       tar xvf "${dl_dir}/${img}"     ;;
			*.tar.*)     tar xvjf "${dl_dir}/${img}"    ;;
			*.tbz2)      tar xvjf "${dl_dir}/${img}"    ;;
			*.tgz)       tar xvzf "${dl_dir}/${img}"    ;;
			*.lzma)      unlzma "${dl_dir}/${img}"      ;;
			*.bz2)       bunzip2 "${dl_dir}/${img}"     ;;
			*.rar)       unrar x -ad "${dl_dir}/${img}" ;;
			*.gz)        gunzip "${dl_dir}/${img}"      ;;
			*.zip)       unzip "${dl_dir}/${img}"       ;;
			*.Z)         uncompress "${dl_dir}/${img}"  ;;
			*.7z)        7z x "${dl_dir}/${img}"        ;;
			*.xz)        unxz "${dl_dir}/${img}"        ;;
		esac
	fi

	if [[ $(file -b --mime-type "${dl_dir}/${img%.*}") == "application/octet-stream" ]]; then
		echo "Preparing image file..."
		sudo guestmount -a "${dl_dir}/${img%.*}" -i "${build_dir}/switchroot/install"
	fi

	echo "Finishing rootfs preparation..."
	mv "${build_dir}/${hekate_bin}" "${build_dir}/lib/firmware/reboot_payload.bin"
	echo -e "/dev/mmcblk0p1	/boot	vfat	rw,relatime	0	2\n" >> "${build_dir}/etc/fstab"
	# sed -r -i 's/^HOOKS=((.*))$/HOOKS=(\1 resize-rootfs)/' "${build_dir}/etc/mkinitcpio.conf"
}

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
docker run --rm -ti -v $(dirname "$(readlink -fm $0)"):/root/ jet-factory:latest go run factory.go -prepare $@
docker run --rm -ti -v $(dirname "$(readlink -fm $0)"):/root/ jet-factory:latest go run factory.go -configs $@
docker run --rm -ti -v $(dirname "$(readlink -fm $0)"):/root/ jet-factory:latest go run factory.go -packages $@
docker run --rm -ti -v $(dirname "$(readlink -fm $0)"):/root/ jet-factory:latest go run factory.go -image $@