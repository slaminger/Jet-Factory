#!/bin/bash
img_url=$1
build_dir="$(readlink -fm $2)"
dl_dir="${build_dir}/downloadedFiles/"
img="${img_url##*/}"

# Hekate files
hekate_version=5.2.0
nyx_version=0.9.0
hekate_url=https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip
hekate_zip=${hekate_url##*/}
hekate_bin=hekate_ctcaer_${hekate_version}.bin

mkdir -p ${build_dir}/{bootloader,switchroot/{install,arch}}

# cd into download directory
mkdir -p ${dl_dir} && cd ${dl_dir}

echo "Downloading Hekate..."
wget -nc -q --show-progress ${hekate_url}

if [ ! -f ${dl_dir}/${img%.*} ]; then
	# Download file if it doesn't exist, or is forced to download.
	echo "Downloading image file..."
	wget -nc --show-progress ${img_url}
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

cd ${dl_dir}

if [ $(file -b --mime-type "${img%.*}") == "application/octet-stream" ]; then
	echo "Preparing image file..."
	loop=$(sudo kpartx -l "${img%.*}" | grep -o -E 'loop[[:digit:]]' | head -1)
	sudo kpartx -a "${img%.*}"
	
	echo "Searching for LVM2 partition type..."
	if [ $(file -b "${img%.*}" | grep "[[:digit:]] : ID=0x8e.*") ]; then
		echo "Found LVM2 partition..."
		rootname=$(sudo lvs | sed 's/root//' | tail -1 | grep -o -E '[[:alpha:]]{3}+')

		echo "Detaching previous LVM2 partition..."
		sudo vgchange -an ${rootname} && sudo vgchange -ay ${rootname}
		sudo mount /dev/mapper/${rootname}-root "${build_dir}/switchroot/install"
	else
		# TODO : Shouldn't try to mount 1st ext2,3,4 partition but biggest
		echo "Found ext2,3,4 partition..."
		num=$(file -b ${img%.*} | grep -o "[[:digit:]] : ID=0x83.*" | cut -d' ' -f1)
		sudo mount /dev/${loop}p${num} "${build_dir}/switchroot/install"
	fi

	echo "Copying files to build directory..."
	cp -prd ${build_dir}/switchroot/install/* ${build_dir} &&
	
	echo "Unmounting partition..."
	sudo umount "${build_dir}/switchroot/install" 
	[[ ! -z ${rootname} ]] && sudo vgchange -an ${rootname}
	sudo kpartx -d ${img%.*}
fi

echo "Finishing rootfs preparation..."
mv "${build_dir}/${hekate_bin}" "${build_dir}/lib/firmware/reboot_payload.bin"
echo -e "/dev/mmcblk0p1	/boot	vfat	rw,relatime	0	2\n" >> "${build_dir}/etc/fstab"
# sed -r -i 's/^HOOKS=((.*))$/HOOKS=(\1 resize-rootfs)/' "${build_dir}/etc/mkinitcpio.conf"
