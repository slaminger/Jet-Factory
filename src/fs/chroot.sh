#!/usr/bin/bash
# CHROOT.SH : Chroot and add qemu binary if necessary
if [[ ${AARCH} != "" ]]; then
	cp "/usr/bin/qemu-${AARCH}-static" "${out}/${NAME}/usr/bin/"
fi

# Mount bind chroot dir
mount --bind "${out}/${NAME}" "${out}/${NAME}"

# Mount dev, proc, sys
mount -t proc proc "${out}/${NAME}/proc/"
mount --rbind /sys "${out}/${NAME}/sys/"
mount --rbind /dev "${out}/${NAME}/dev/"

# Copy vuild script
cp "$(dirname ${cwd})/configs/examples/${CHROOT_SCRIPT}" "${out}/${NAME}"

# Actual chroot
arch-chroot "${out}/${NAME}" /bin/bash /${CHROOT_SCRIPT}

# Remove build script
rm "${out}/${NAME}/${CHROOT_SCRIPT}"

# Umount dev, proc, sys
umount -R ${out}/${NAME}/{dev,proc,sys}/

# Umount chroot dir
umount "${out}/${NAME}"

if [[ ${AARCH} != "" ]]; then
	rm "${out}/${NAME}/usr/bin/qemu-${AARCH}-static"
fi