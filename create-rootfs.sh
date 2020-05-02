#!/usr/bin/bash

build_dir=/root/l4t/
url=ROOTFS_URL
archive="$(cut -d \"/\" -f1 ${url})"
hekate_version=5.2.0
nyx_version=0.9.0

prepareFiles() {
	mkdir -p ${build_dir}/{{r,b}ootfs,tmp,switchroot/install/}
	
	wget https://github.com/CTCaer/hekate/releases/download/v${hekate_version}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip -P ${build_dir}
	wget ${url} -P ${build_dir} &&

	# TODO : Extract and copy rootfs files condition
	if ![[ ${archive} == *.raw.*]] || ; then
		bsdtar xpf ${build_dir}/${archive} -C ${build_dir}/rootfs/
	elif [[ ${archive} == *.raw.xz ]]; then
	fi

	unzip ${build_dir}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip hekate_ctcaer_${hekate_version}.bin
	mv hekate_ctcaer_${hekate_version}.bin ${build_dir}/rootfs/lib/firmware/reboot_payload.bin
	rm ${build_dir}/hekate_ctcaer_${hekate_version}_Nyx_${nyx_version}.zip
}

setup() {
	echo -e "/dev/mmcblk0p1	/boot	vfat	rw,relatime	0	2\n" >> ${build_dir}/rootfs/etc/fstab
	sed -i 's/^HOOKS=(\(.*\))$/HOOKS=(\1 resize-rootfs)/' ${build_dir}/rootfs/etc/mkinitcpio.conf
	cp /usr/bin/qemu-aarch64-static ${build_dir}/rootfs/usr/bin/

	mount --bind ${build_dir}/rootfs ${build_dir}/rootfs
	mount --bind ${build_dir}/bootfs ${build_dir}/rootfs/boot/

	if ${arch}; then
		# Workaround for flakiness of `pt` mirror.
		sed -i 's/mirror.archlinuxarm.org/de.mirror.archlinuxarm.org/g' ${build_dir}/rootfs/etc/pacman.d/mirrorlist
		echo "[switch]\nSigLevel = Optional\nServer = https://9net.org/l4t-arch/" >> ${build_dir}/rootfs/etc/pacman.conf

		# Install Packages
		arch-chroot ${build_dir}/rootfs/
		# Install configs
		arch-chroot ${build_dir}/rootfs/

		rm ${build_dir}/rootfs/etc/pacman.d/gnupg/S.gpg-agent*
	elif ${fedora}; then
		# Install Packages
		arch-chroot ${build_dir}/rootfs/
		# Install configs
		arch-chroot ${build_dir}/rootfs/
	elif ${opensuse}; then
		# Install Packages
		arch-chroot ${build_dir}/rootfs/
		# Install configs
		arch-chroot ${build_dir}/rootfs/
	fi

	rm ${build_dir}/rootfs/usr/bin/qemu-aarch64-static
	umount -R ${build_dir}/rootfs/{,boot/}
	mv ${build_dir}/bootfs/* ${build_dir}/
}

buildimg() {
	size=$(du -hs -BM ${build_dir}/rootfs/ | head -n1 | awk '{print int($1/4)*4 + 4 + 512;}')M
	echo "Estimated rootfs size: $size"

	dd if=/dev/zero of=${build_dir}/switchroot/install/l4t.img bs=1 count=0 seek=$size
	
	loop=`losetup --find`
	losetup ${loop} ${build_dir}/switchroot/install/l4t.img

	mkfs.ext4 ${loop}
	mount ${loop} ${build_dir}/tmp

	mv ${build_dir}/rootfs/* ${build_dir}/tmp/

	umount ${loop}
	losetup -d ${loop}

	cd ${build_dir}/switchroot/install/
	split -b4290772992 --numeric-suffixes=0 l4t-fedora.img l4t.

	rm -rf ${build_dir}/{{b,r}ootfs/,tmp/,switchroot/install/l4t.img}
}

echo "\nPreparing required files\n"
prepareFiles &&
echo "\nCreating rootfs\n"
setup &&
echo "\nBuilding image file\n"
buildimg &&
echo "Done!\n"