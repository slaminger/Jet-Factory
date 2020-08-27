#!/bin/bash
# CHROOT.SH : Arch-Chroot with qemu binary (if needed)
if [ ! -f /proc/sys/fs/binfmt_misc/register ]; then
	if ! mount binfmt_misc -t binfmt_misc /proc/sys/fs/binfmt_misc; then
        exit 1
    fi
fi

if [[ -n "${AARCH}" && ! -e "/proc/sys/fs/binfmt_misc/qemu-${AARCH}" ]]; then
	wget -L -q -nc --show-progress https://raw.githubusercontent.com/dbhi/qus/main/register.sh
	chmod +x register.sh
	./register.sh -s -- -p "${AARCH}"
	cp "/usr/bin/qemu-${AARCH}-static" "${out}/${NAME}/usr/bin/"
fi

# Mount bind chroot dir
mount --bind "${out}/${NAME}" "${out}/${NAME}"

# Copy build script
cp "$(dirname "${cwd}")/configs/${DEVICE}/files/${CHROOT_SCRIPT}" "${out}/${NAME}"

# Handle resolv.conf
[[ -e "${out}/${NAME}/etc/resolv.conf" && ! -L "${out}/${NAME}/etc/resolv.conf" ]] && \
	cp "${out}/${NAME}/etc/resolv.conf" "${out}/${NAME}/etc/resolv.conf.bak"

echo "namserver 8.8.8.8" > resolv.conf
cp resolv.conf "${out}/${NAME}/etc/resolv.conf"

# Actual chroot
arch-chroot "${out}/${NAME}" /bin/bash /"${CHROOT_SCRIPT}"

# Clean temp files
rm "${out}/${NAME}/${CHROOT_SCRIPT}" resolv.conf

[[ -e "${out}/${NAME}/etc/resolv.conf.bak" ]] && \
	cp "${out}/${NAME}/etc/resolv.conf.bak" "${out}/${NAME}/etc/resolv.conf"

# Umount chroot dir
umount "${out}/${NAME}"

if [[ -n "${AARCH}" && -e "/proc/sys/fs/binfmt_misc/qemu-${AARCH}" ]]; then
	./register.sh -- -r
	rm "${out}/${NAME}/usr/bin/qemu-${AARCH}-static"
	rm ./register.sh
fi
