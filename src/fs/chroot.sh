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

# Add cache dir configuration
if [ ! -z "$CACHE_DIR" ]; then
  mkdir "${out}/cache" &> /dev/null || true
  mount --bind "${out}/cache" "${out}/${NAME}/${CACHE_DIR}" || exit
fi

# Copy build script
cp "$(dirname "${cwd}")/configs/${DEVICE}/files/${CHROOT_SCRIPT}" "${out}/${NAME}"

# Handle resolv.conf
if [[ -e "${out}/${NAME}/etc/resolv.conf" ]]; then
	cp "${out}/${NAME}/etc/resolv.conf" "${out}/${NAME}/etc/resolv.conf.bak"
	echo "namserver 8.8.8.8" > resolv.conf
	cp resolv.conf "${out}/${NAME}/etc/resolv.conf"
fi

# Actual chroot
arch-chroot "${out}/${NAME}" /bin/bash /"${CHROOT_SCRIPT}"

# Clean temp files
rm -rf "${out}/${NAME}/${CHROOT_SCRIPT}" "${out}/cache"

if [[ -e "${out}/${NAME}/etc/resolv.conf.bak" ]]; then
	cp "${out}/${NAME}/etc/resolv.conf.bak" "${out}/${NAME}/etc/resolv.conf"
	rm -rf resolv.conf
fi

# Umount chroot dir
umount "${out}/${NAME}"

if [[ -n "${AARCH}" && -e "/proc/sys/fs/binfmt_misc/qemu-${AARCH}" ]]; then
	./register.sh -- -r
	rm -rf "${out}/${NAME}/usr/bin/qemu-${AARCH}-static"
	rm -rf ./register.sh
fi
