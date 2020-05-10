#!/bin/bash
build_dir=$1
img=$2

mkdir -p ${build_dir}

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

if [[ ${hekate} == true ]]; then
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