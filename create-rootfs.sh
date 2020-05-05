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
cwd="$(dirname "$(readlink -f "$0")")"
build_dir="${cwd}/build"
dl_dir="${cwd}/dl"

# Distro specific variables
selection="$(echo ${@: -1} | tr '[:upper:]' '[:lower:]')"
img_url="$(head -1 ${cwd}/install/${selection}/urls)"
img="${img_url##*/}"
img_sig_url="${img_url}.md5"
img_sig="${img_sig_url##*/}"

# Hekate files
hekate_version=5.2.0
nyx_version=0.9.0
hekate_url=https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip
hekate_zip=${hekate_url##*/}
hekate_bin=hekate_ctcaer_${hekate_version}.bin

(
	if [[ ${selection} == "arch" ]]; then
		img_sig_url="${img_url}.md5"
		validate_command="md5sum --status -c "${img_sig}""
	elif [[ ${selection} == "fedora" ]]; then
		img_sig_url=https://download.fedoraproject.org/pub/fedora/linux/releases/31/Server/aarch64/images/Fedora-Server-31-1.9-aarch64-CHECKSUM
		validate_command="sha256sum --status -c "${img_sig}""
		img_sig_url="Fedora-Server-31-1.9-aarch64-CHECKSUM"
	else
		echo "$0: invalid distro option: $1"
		usage
		exit 1
	fi
)

usage() {
    echo "Usage: $0 [options] <distribution-name>"
    echo "Options:"
	echo " -d, --docker          	Build with Docker"
	echo " --hekate                 Build for Hekate"
	echo " -k, --keep          		Keep downloaded files"
    echo " -s, --staging            Install built local packages"
    echo " -h, --help               Show this help text"
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
	# Will check for a dir with the name of the arcihve without extensions
	[[ "${dl_dir}/${img%.*}" =~ '.tar' ]] && cd ${build_dir} &&
	dtrx -f ${dl_dir}/${img} && tar xf ${img%.*} && rm ${img%.*}

	# Will check for a dir with the name of the arcihve with 1st extension only
	[[ ! -f "${dl_dir}/${img%.*}" ]] &&
	dtrx -f ${dl_dir}/${img}

	if [[ $(file -b --mime-type "${dl_dir}/${img%.*}") == "application/octet-stream" ]]; then
		echo "Searching for LVM2 partition type..."
		kpartx -av "${dl_dir}/${img%.*}"

		# Use this for non LVM partition
		if [[ $(file -b "${dl_dir}/${img%.*}" | grep "[[:digit:]] : ID=0x8e.*") ]]; then
			echo "Found LVM2 partition..."
			vgchange -ay
			sleep 2
			mount_part="/dev/$(vgchange -ay | grep -o '".*"' | sed 's/"//g')/root"
			sleep 2
		else
			echo "Didn't found LVM2 partition..." && $loop
			mount_part="/dev/mapper/${loop##*/}p$(file -b "${dl_dir}/${img%.*}" |  grep -o "partition 2.*" | grep -o "[[:digit:]] : ID=0x83.*" | cut -d' ' -f1)"
			sleep 2
		fi
		
		echo "Mounting partition..."
		mount ${mount_part} "${build_dir}/bootloader/"

		echo "Copying files to build directory..."
		cp -prd "${build_dir}/bootloader/" ${build_dir} 2>/dev/null
		
		[[ $(file -b "${dl_dir}/${img%.*}" | grep "[[:digit:]] : ID=0x8e.*") ]] &&
		vgchange -an
		umount "${build_dir}/bootloader/"
		kpartx -d "${dl_dir}/${img%.*}"
		losetup -d ${loop}
	fi

	cd ${build_dir}
	echo "Extracting Hekate..."
	dtrx -f ${dl_dir}/${hekate_zip}
	cd ${cwd}
}

Main() {
	echo "Create required directories..."
	mkdir -p ${dl_dir} ${build_dir}/{switchroot/install/pkgs,bootloader/}

	echo "Downloading image..."
	GetImgFiles
	
	echo "Extracting image..."
	ExtractFiles

	echo "Copying files to rootfs..."
	[[ ${staging} == true ]] && cp -r "${cwd}/install/${selection}/*/*/${pkg_types}" "${build_dir}/pkgs/" 2>/dev/null
	cp ${cwd}/install/${selection}/{build-stage2.sh,base-pkgs} ${build_dir}
	mv "${build_dir}/${hekate_bin}" ${build_dir}/lib/firmware/reboot_payload.bin
	
	echo "Pre chroot setup..."
	echo -e "/dev/mmcblk0p1	/boot	vfat	rw,relatime	0	2\n" >> ${build_dir}/etc/fstab
	sed -r -i 's/^HOOKS=((.*))$/HOOKS=(\1 resize-rootfs)/' ${build_dir}/etc/mkinitcpio.conf
	chmod +x ${build_dir}/build-stage2.sh
	
	mount --bind ${build_dir} ${build_dir} &&
	mount --bind  "${build_dir}/bootloader/" "${build_dir}/boot/"
	
	echo "Chrooting..."
	arch-chroot ${build_dir} ./build-stage2.sh || exit 1

	echo "Post chroot cleaning..."
	umount "${build_dir}/boot/"
	umount ${build_dir}
	rm -rf ${build_dir}/{base-pkgs,build-stage2.sh,pkgs/,usr/bin/qemu-aarch64-static}

	echo "Creating final ${format} partition..."
	size=$(du -hs -BM ${build_dir} | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M
	dd if=/dev/zero of="${img%%.*}" bs=1 count=0 seek=${size} status=noxfer
	
	echo "Creating"${img%%.*}" with ${format} format..."
	yes | mkfs.${format} "${img%%.*}"
	mount -o loop "${img%%.*}" "${build_dir}/switchroot/install/"
	cp -prd ${build_dir}/* "${build_dir}/switchroot/install/" 2>/dev/null
	umount ${loop}

	if [[ ${hekate} == true ]]; then
		echo "Creating Hekate installable files..."
		split -b4290772992 --numeric-suffixes=0 "${img%%.*}" l4t.
	
		echo "Compressing hekate folder..."
		7z a ${cwd}/"SWR-${img%%.*}.7z" ${build_dir}/{bootloader,switchroot}
	else
		echo "Creating ${img%%.*}.fat32 with ${format} format..."
		size=$(du -hs -BM "${build_dir}/bootloader" | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M
		dd if=/dev/zero of="${img%%.*}.fat32" bs=1 count=0 seek=${size} status=noxfer
		losetup ${loop} "${img%%.*}.fat32"
		yes | mkfs.vfat -F 32 ${loop}
		
		mount -o loop "${img%%.*}.fat32" "${build_dir}/boot"
		cp -r ${build_dir}/bootloader/* "${build_dir}/boot" 2>/dev/null
		umount ${loop}
	
		echo "Creating final image: ${img%%.*}.img..."
		dd if="${img%%.*}.fat32" bs=1M count=99 skip=1 of="SWR-${img%%.*}.img" status=noxfer
		dd if="${img%%.*}" bs=1M count=10 of="SWR-${img%%.*}.img" oflag=append conv=notrunc status=noxfer
	fi

	echo "Cleaning up files..."
	losetup -d ${loop}
	umount ${build_dir}
	rm -rf ${build_dir} ${cwd}/${img%%.*}{,.fat32}
	echo "Done!"
}

# Parse arguments
options=$(getopt -n $0 -o dhks --long docker,keep,staging,hekate,help -- "$@")

# Check for errors in arguments or if no name was provided
if [[ $? != "0" ]] || [[ options =~ "${@: -1}" ]]; then usage; exit 1; fi

# Evaluate arguments
eval set -- "$options"
while true; do
    case "$1" in
	-d | --docker) docker=true; shift ;;
	--hekate) hekate=true; shift ;;
	-k | --keep) keep=true; shift ;;
    -s | --staging) staging=true; shift ;;
    ? | -h | --help) usage; exit 0 ;;
    -- ) shift; break ;;
    esac
done

echo "Cleaning up old Files..."
rm -rf ${build_dir} ${cwd}/${img%%.*}{,.fat32}
[[ ${keep} == true ]] && rm -rf ${dl_dir}

if [[ ${docker} == true ]]; then
	echo "Starting Docker..."
	systemctl start docker.{socket,service}
	echo "Building Docker image..."
	docker image build -t l4t-builder:1.0 .
	echo "Running container..."
	docker run --privileged --cap-add=SYS_ADMIN --rm -it -v ${cwd}:/builder l4t-builder:1.0 /bin/bash /builder/create-rootfs.sh $(echo "$options" | sed -E 's/-(d|-docker)//g' | grep -o -E '\-\-[[:alpha:]]+' | tr '\r\n' ' ') ${selection}
	exit 0
fi

Main && exit 0