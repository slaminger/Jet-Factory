#!/bin/bash
# CHROOT.SH : Arch-Chroot with qemu binary (if needed)
if [[ $(uname -m) != ${AARCH} ]]; then
	if [ ! -f /proc/sys/fs/binfmt_misc/register ]; then
		if ! mount binfmt_misc -t binfmt_misc /proc/sys/fs/binfmt_misc; then
	        exit 1
	    fi
	fi

	if [[ -n "${AARCH}" && ! -e "/proc/sys/fs/binfmt_misc/qemu-${AARCH}" ]]; then
		wget -L -q -nc --show-progress https://raw.githubusercontent.com/dbhi/qus/main/register.sh -P "${out}/downloadedFiles/"
		chmod +x "${out}/downloadedFiles/register.sh"
		"${out}/downloadedFiles/register.sh" -s -- -p "${AARCH}"
		cp "/usr/bin/qemu-${AARCH}-static" "${out}/${NAME}/usr/bin/"
	fi
fi

# Mount bind chroot dir
mount --bind "${out}/${NAME}" "${out}/${NAME}"

# Add cache dir configuration
if [ -n "$CACHE_DIR" ]; then
  mkdir "${out}/cache" &> /dev/null || true
  mount --bind "${out}/cache" "${out}/${NAME}/${CACHE_DIR}" || exit
fi

# Copy build script
cp "$(dirname "${cwd}")/configs/${DEVICE}/files/${CHROOT_SCRIPT}" "${out}/${NAME}"

# Handle resolv.conf
cp --dereference /etc/resolv.conf "${out}/${NAME}/etc/resolv.conf"

# Actual chroot
arch-chroot "${out}/${NAME}" /bin/bash /"${CHROOT_SCRIPT}"

# Clean temp files
rm -rf "${out}/${NAME}/${CHROOT_SCRIPT}" "${out}/cache"

# Umount chroot dir
umount "${out}/${NAME}"

if [[ -n "${AARCH}" && ! -e "/proc/sys/fs/binfmt_misc/qemu-${AARCH}" && $(uname -m) != ${AARCH} ]]; then
	"${out}/downloadedFiles/register.sh" -- -r
	rm -rf "${out}/${NAME}/usr/bin/qemu-${AARCH}-static" "${out}/downloadedFiles/register.sh"
fi
