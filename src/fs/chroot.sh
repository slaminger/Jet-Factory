#!/bin/bash
# CHROOT.SH : Chroot and add qemu binary if necessary
if [[ ${AARCH} != "" ]]; then
	wget -L -q -nc --show-progress https://raw.githubusercontent.com/dbhi/qus/main/register.sh
	chmod +x register.sh
	./register.sh -s -- -p "${AARCH}"
	cp "/usr/bin/qemu-${AARCH}-static" "${out}/${NAME}/usr/bin/"
fi

# Mount bind chroot dir
mount --bind "${out}/${NAME}" "${out}/${NAME}"

# Copy build script
cp "$(dirname "${cwd}")/configs/examples/${CHROOT_SCRIPT}" "${out}/${NAME}"

echo "namserver 8.8.8.8" > resolv.conf
[[ -e "${out}/${NAME}/etc/resolv.conf" && ! -L "${out}/${NAME}/etc/resolv.conf" ]] && cp "${out}/${NAME}/etc/resolv.conf" "${out}/${name}/etc/resolv.conf.bak"
cp resolv.conf "${out}/${NAME}/etc/resolv.conf"

# Actual chroot
arch-chroot "${out}/${NAME}" /bin/bash /"${CHROOT_SCRIPT}"

# Remove build script
rm "${out}/${NAME}/${CHROOT_SCRIPT}"

[[ -e "${out}/${NAME}/etc/resolv.conf.bak" ]] && cp "${out}/${NAME}/etc/resolv.conf.bak" "${out}/${name}/etc/resolv.conf"
rm resolv.conf

# Umount chroot dir
umount "${out}/${NAME}"

if [[ ${AARCH} != "" ]]; then
	./register.sh -- -r
	rm "${out}/${NAME}/usr/bin/qemu-${AARCH}-static"
	rm ./register.sh
fi
