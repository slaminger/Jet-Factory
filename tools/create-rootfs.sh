#!/bin/env bash

# Setup variables
docker=false
staging=false
hekate=false
keep=false

pkg_types=*.{pkg.*,rpm,deb}
format=ext4
loop=`losetup --find`

# Folders
parent="$(dirname "$(dirname "$(readlink -fm "$0")")")"
cwd="$(dirname "$(readlink -f "$0")")"
dl_dir="${cwd}/builds/dl"

# Distro specific variables
selection="$(echo ${@: -1} | tr '[:upper:]' '[:lower:]')"
build_dir="${cwd}/builds/${selection}-build"
img_url="$(head -1 ${parent}/configs/${selection}/urls)"
img="${img_url##*/}"

# Hekate files
hekate_version=5.2.0
nyx_version=0.9.0
hekate_url=https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip
hekate_zip=${hekate_url##*/}
hekate_bin=hekate_ctcaer_${hekate_version}.bin

Usage() {
    echo "Usage: $0 [options] <distribution-name>"
    echo "Options:"
	echo " -d, --docker          	Build with Docker"
	echo " --hekate                 Build for Hekate"
	echo " -k, --keep          		Keep downloaded files"
    echo " -s, --staging            Install built local packages"
    echo " -h, --help               Show this help text"
}

# TODO : Replace by parsed YAML
SetDistro() {
	if [[ ${selection} == "arch" ]]; then
		img_sig_url="${img_url}.md5"
		img_sig="${img_sig_url##*/}"
		validate_command="md5sum --status -c "${dl_dir}/${img_sig}""
	elif [[ ${selection} == "fedora" ]]; then
		img_sig_url=https://download.fedoraproject.org/pub/fedora/linux/releases/31/Server/aarch64/images/Fedora-Server-31-1.9-aarch64-CHECKSUM
		img_sig="${img_sig_url##*/}"
		validate_command="sha256sum --status -c "${dl_dir}/${img_sig}""
	else
		echo "$0: invalid distro option: $1"
		Usage
		exit 1
	fi
}

GetImgFiles() {
	echo "Downloading Hekate..."
	wget -nc -P ${dl_dir} -q --show-progress ${hekate_url} -O ${dl_dir}/${hekate_zip}
	
	# Check image against signature
	echo "Validating image..."
	$validate_command
	[[ $? != "0" ]] && echo "Signature check passed!"
	
	# cd into download directory
	cd ${dl_dir}

	# Download file if it doesn't exist, or is forced to download.
	echo "Downloading image file..."
	wget -nc -q --show-progress ${img_url} -O "${dl_dir}/${img}"
	
	# Download signature file
	echo "Downloading signature file..."
	wget -q --show-progress ${img_sig_url} -O "${dl_dir}/${img_sig}"
}

ExtractFiles() {
	cd ${build_dir}
	echo "Extracting downloaded archive file..."
	[[ ! -f "${dl_dir}/${img%.*}" ]] && dtrx -f ${dl_dir}/${img}

	[[ "${dl_dir}/${img%.*}" =~ '.tar' ]] && cd ${build_dir} &&
	dtrx -f ${dl_dir}/${img} && tar xf ${img%.*} &>/dev/null &&
	rm ${img%.*}

	echo "Searching for image file..."
	if [[ $(file -b --mime-type "${dl_dir}/${img%.*}") == "application/octet-stream" ]]; then
		
		echo "Preparing image file..."
		loop=$(kpartx -l "${dl_dir}/${img%.*}" | grep -o -E 'loop[[:digit:]]' | head -1)
		kpartx -a "${dl_dir}/${img%.*}"
		
		echo "Searching for LVM2 partition type..."
		if [[ $(file -b "${dl_dir}/${img%.*}" | grep "[[:digit:]] : ID=0x8e.*") ]]; then

			echo "Found LVM2 partition..."  && echo "Searching for rootfs partition..."
			rootname=$(lvs | sed 's/root//' | tail -1 | grep -o -E '[[:alpha:]]{3}+')

			echo "Detaching previous LVM2 partition..."
			vgchange -an ${rootname} && vgchange -ay ${rootname}
			mount /dev/mapper/${rootname}-root "${build_dir}/switchroot/install"
		else
			# TODO : Shouldn't try to mount 1st ext2,3,4 partition but biggest
			echo "Found ext2,3,4 partition..."
			num=$(file -b "${dl_dir}/${img%.*}" | grep -o "[[:digit:]] : ID=0x83.*" | cut -d' ' -f1)
			mount /dev/${loop}p${num} "${build_dir}/switchroot/install"
		fi

		echo "Copying files to build directory..."
		cp -prd ${build_dir}/switchroot/install/* ${build_dir} &&
		
		echo "Unmounting partition..."
		umount "${build_dir}/switchroot/install" 
		[[ ! -z ${rootname} ]] && vgchange -an ${rootname}
		kpartx -d "${dl_dir}/${img%.*}"
	fi

	echo "Extracting Hekate..."
	dtrx -f ${dl_dir}/${hekate_zip}
	cd ${cwd}
}

PreChroot() {
	[[ ${staging} == true ]] && echo "Copying staging files to rootfs..." &&
	cp -r ${parent}/configs/${selection}/*.${pkg_types} "${build_dir}/pkgs/" 2>/dev/null
	
	echo "Copying files to rootfs..."
	cp /usr/bin/qemu-aarch64-static ${build_dir}/usr/bin/
	cp ${parent}/configs/${selection}/{build-stage2.sh,base-pkgs} ${build_dir}
	mv "${build_dir}/${hekate_bin}" "${build_dir}/lib/firmware/reboot_payload.bin"
	
	echo "Pre chroot setup..."
	echo -e "/dev/mmcblk0p1	/boot	vfat	rw,relatime	0	2\n" >> "${build_dir}/etc/fstab"
	# sed -r -i 's/^HOOKS=((.*))$/HOOKS=(\1 resize-rootfs)/' "${build_dir}/etc/mkinitcpio.conf"
	chmod +x "${build_dir}/build-stage2.sh"
}

Chroot() {
	mount --bind ${build_dir} ${build_dir} &&
	mount --bind  "${build_dir}/bootloader/" "${build_dir}/boot/"
	
	echo "Chrooting..."
	# TODO : Golang : Build stage2 will be replaced by configs and packages installation
	arch-chroot ${build_dir} /bin/bash /build-stage2.sh || exit 1

	echo "Post chroot cleaning..."
	umount "${build_dir}/boot/" && umount ${build_dir}
	rm -rf ${build_dir}/{base-pkgs,build-stage2.sh,pkgs/,usr/bin/qemu-aarch64-static}
}

PostChroot() {
	size=$(du -hs -BM ${build_dir} | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M
	echo "Creating final ${format} partition... Estimated size : ${size}"
	dd if=/dev/zero of="${cwd}/${img%%.*}" bs=1 count=0 seek=${size} status=noxfer
	
	echo "Creating"${img%%.*}" with ${format} format..."
	ext4mnt=$(losetup --partscan --find --show "${cwd}/${img%%.*}")
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
		split -b4290772992 --numeric-suffixes=0 ${cwd}/"${img%%.*}" l4t.
		
		echo "Compressing hekate folder..."
		if [[ "${dl_dir}/${img%.*}" =~ '.tar' ]]; then 
			7z a ${cwd}/"SWR-${img%%.*}.7z" ${build_dir}/{bootloader,switchroot}
		else
			7z a ${cwd}/"SWR-${img%.*}.7z" ${build_dir}/{bootloader,switchroot}
		fi
	else
		echo "Creating ${img%%.*}.fat32 with ${format} format..."
		size=$(du -hs -BM "${build_dir}/bootloader" | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M
		dd if=/dev/zero of="${img%%.*}.fat32" bs=1 count=0 seek=${size} status=noxfer
		fat32mnt=$(losetup --partscan --find --show "${img%%.*}".fat32)
		yes | mkfs.vfat -F 32 ${fat32mnt}
		
		echo "Copying files to fat32 partition..."
		mount ${fat32mnt} "${build_dir}/boot"
		cp -r ${build_dir}/bootloader/* "${build_dir}/boot" 2>/dev/null
		umount ${fat32mnt} && losetup -d ${fat32mnt}
	
		echo "Creating final image: ${img%%.*}.img..."
		dd if="${img%%.*}.fat32" bs=1M count=99 skip=1 of="SWR-${img%%.*}.img" status=noxfer
		dd if="${img%%.*}" bs=1M count=10 of="SWR-${img%%.*}.img" oflag=append conv=notrunc status=noxfer
	fi
}

BuildEngine() {
	if [ ! -f /.dockerenv ]; then
		echo "Cleaning up old builds..."
		rm -rf ${build_dir} ${cwd}/${img%%.*}{,.fat32}
		[[ ${keep} != true ]] && echo "Keeping previously downloaded files..." && rm -rf ${dl_dir}

		echo "Create required directories..."
		mkdir -p ${dl_dir} ${build_dir}/{switchroot/install/,bootloader,pkgs}

		echo "Setting distro parameters..."
		SetDistro

		echo "Downloading image..."
		GetImgFiles
		
		echo "Extracting image..."
		cd ${build_dir}
		ExtractFiles
	fi

	if [[ ${docker} == true ]]; then
		echo "Starting Docker..."
		systemctl start docker.{socket,service}

		echo "Building Docker image..."
		docker image build -t l4t-builder:1.0 .

		echo "Running container..."
		docker run --privileged --cap-add=SYS_ADMIN --rm -it -v ${parent}:/root l4t-builder:1.0 /bin/bash /root/helpers/create-rootfs.sh $(echo "$options" | sed -E 's/-(d|-docker)//g' | grep -o -E '\-+[[:alpha:]]+' | tr '\r\n' ' ') ${selection}
		exit 0
	fi

	echo "Pre Chroot setup..."
	PreChroot

	echo "Chrooting..."
	Chroot

	echo "Post Chroot setup..."
	PostChroot

	echo "Cleaning up files..."
	rm -rf ${build_dir} ${cwd}/${img%%.*}{,.fat32}
	echo "Done!"
}

# Parse arguments
options=$(getopt -n $0 -o dhks --long docker,keep,staging,hekate,help -- "$@")

# Check for errors in arguments or if no name was provided
if [[ $? != "0" ]] || [[ options =~ "${@: -1}" ]]; then Usage; exit 1; fi
[[ `whoami` != root ]] && echo "Run this as root!" && exit 1

# Evaluate arguments
eval set -- "$options"
while true; do
    case "$1" in
	-d | --docker) docker=true; shift ;;
	--hekate) hekate=true; shift ;;
	-k | --keep) keep=true; shift ;;
    -s | --staging) staging=true; shift ;;
    ? | -h | --help) Usage; exit 0 ;;
    -- ) shift; break ;;
    esac
done

BuildEngine && exit 0