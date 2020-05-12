#!/bin/bash
img_url=${@: -2}
build_dir="$(readlink -fm ${@: -1})"
dl_dir="${build_dir}/downloadedFiles/"
img="${img_url##*/}"

# Hekate files
hekate_version=5.2.0
nyx_version=0.9.0
hekate_url=https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip
hekate_zip=${hekate_url##*/}
hekate_bin=hekate_ctcaer_${hekate_version}.bin

PreChroot() {
	echo "Pre chroot setup..."
	cp /usr/bin/qemu-aarch64-static ${build_dir}/usr/bin/
	mount --bind ${build_dir} ${build_dir} &&
	mount --bind  "${build_dir}/bootloader/" "${build_dir}/boot/"
}

PostChroot() {
	echo "Post chroot cleaning..."
	umount "${build_dir}/boot/" && umount ${build_dir}
	rm ${build_dir}/usr/bin/qemu-aarch64-static
}

PrepareFiles() {
	mkdir -p ${build_dir}/{bootloader,switchroot/{install,arch}} ${dl_dir}
	cd ${dl_dir}
	
	wget -nc -q --show-progress ${hekate_url}
	unzip ${hekate_zip} -d ${build_dir}
	wget -nc --show-progress ${img_url}		

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

	echo "Copying files to build directory..."
	cp -prd ${build_dir}/switchroot/install/* ${build_dir}

	echo "Unmounting partition..."
	sudo umount "${build_dir}/switchroot/install" 

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
	echo "Creating final ${format} partition... Estimated size : ${size}"
	dd if=/dev/zero of=${img} bs=1 count=0 seek=${size} status=noxfer

	echo "Creating${img} with ${format} format..."
	ext4mnt=$(losetup --partscan --find --show ${img})
	yes | mkfs.${format} ${ext4mnt}

	echo "Copying files to ${format} partition..."
	mount ${ext4mnt} "${build_dir}/switchroot/install/"
	cp -prd ${build_dir}/* "${build_dir}/switchroot/install/" 2>/dev/null
	cp -prd ${build_dir}/boot ${build_dir}/switchroot/install/ 2>/dev/null

	echo "Removing unneeded folders from partiton..."
	rm -rf ${build_dir}/{switchroot/install/,bootloader/{switchroot/,bootloader/},*.reg}
	umount ${ext4mnt} && losetup -d ${ext4mnt}

	if [[ $1 == "hekate" ]]; then
		echo "Creating Hekate installable files..."
		cd "${build_dir}/switchroot/install/"
		split -b4290772992 --numeric-suffixes=0 ${img} l4t.
		
		echo "Compressing hekate folder..."
		if [[ "${dl_dir}/${img%.*}" =~ '.tar' ]]; then 
			7z a "SWR-${img%%.*}.7z" ${build_dir}/{bootloader,switchroot}
		else
			7z a "SWR-${img%.*}.7z" ${build_dir}/{bootloader,switchroot}
		fi
	else
		echo "Creating ${img%%.*}.fat32 with ${format} format..."
		size=$(du -hs -BM "${build_dir}/bootloader" | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M
		dd if=/dev/zero of="${img%%.*}.fat32" bs=1 count=0 seek=${size} status=noxfer
		fat32mnt=$(losetup --partscan --find --show ${img}.fat32)
		yes | mkfs.vfat -F 32 ${fat32mnt}
		
		echo "Copying files to fat32 partition..."
		mount ${fat32mnt} "${build_dir}/boot"
		cp -r ${build_dir}/bootloader/* "${build_dir}/boot" 2>/dev/null
		umount ${fat32mnt} && losetup -d ${fat32mnt}

		echo "Creating final image: ${img%%.*}.img..."
		dd if="${img%%.*}.fat32" bs=1M count=99 skip=1 of="SWR-${img%%.*}.img" status=noxfer
		dd if=${img} bs=1M count=10 of="SWR-${img%%.*}.img" oflag=append conv=notrunc status=noxfer
	fi
}

if [[ $@ =~ "pre" ]]; then
	PreChroot
fi

if [[ $@ =~ "pkg" ]]; then
	findPkgManager
elif [[ $@ =~ "files" ]]; then
	PrepareFiles
elif [[ $@ =~ "image" ]]; then
	CreateImage
fi

if [[ $@ =~ "post" ]]; then
	PostChroot
fi