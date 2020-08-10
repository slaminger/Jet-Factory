#!/bin/bash
# CHROOT.SH : Chroot and add qemu binary if necessary
if [[ ${AARCH} != "" ]]; then
	wget -L https://raw.githubusercontent.com/dbhi/qus/main/register.sh
	chmod +x register.sh
	./register.sh -s -- -p ${AARCH}
	cp "/usr/bin/qemu-${AARCH}-static" "${out}/${NAME}/usr/bin/"
fi

# Mount bind chroot dir
mount --bind "${out}/${NAME}" "${out}/${NAME}"

# Copy build script
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
	./register.sh -- -r
	rm "${out}/${NAME}/usr/bin/qemu-${AARCH}-static"
	rm ./register.sh
fi
