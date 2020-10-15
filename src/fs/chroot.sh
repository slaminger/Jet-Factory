#!/bin/bash
# CHROOT.SH : Arch-Chroot with QEMU CPU architecture emulation support

same_cpu_arch=1
if [[ -n "${AARCH}" ]] && [[ $(uname -m) != ${AARCH} ]]; then
	same_cpu_arch=0
else
	AARCH=$(uname -m)
	echo -e "\n\t\tAssuming target and host use the same CPU architecture: ${AARCH}\n"
fi

# Check if architecture is already registered in the .lock file
lock="$(grep -qxF -- "${AARCH}" "${out}/.lock")"

# If an architecture is  already registered increment the counter
if [[ ${lock} = 1 ]]; then
	# Get current lock count
	lock_count=$(sed 's/'${AARCH}' //g' "${out}/.lock")

	# Increment lock count in lock file
	sed -i 's/'${AARCH}' '${lock_count}'/'${AARCH}' '$((lock_count+1))'/g' "${out}/.lock"

	# Increment lock count variable
	lock_count=$((lock_count+1))
else
	echo "${AARCH} 1" > "${out}/.lock"
fi

if [[ ${same_cpu_arch} = 0 ]] && [[ ${lock} = 1 ]]; then
	if [ ! -f /proc/sys/fs/binfmt_misc/register ]; then
		if ! mount binfmt_misc -t binfmt_misc /proc/sys/fs/binfmt_misc; then
	        exit 1
	    fi
	fi
	
	if [[ ! -e "/proc/sys/fs/binfmt_misc/qemu-${AARCH}" ]]; then
		wget -L -q -nc --show-progress https://raw.githubusercontent.com/dbhi/qus/main/register.sh -P "${out}/downloadedFiles/"
		chmod +x "${out}/downloadedFiles/register.sh"
		"${out}/downloadedFiles/register.sh" -s -- -p "${AARCH}"
		cp "/usr/bin/qemu-${AARCH}-static" "${out}/${NAME}/usr/bin/"
	fi
fi

# Mount bind chroot dir
mount --bind "${out}/${NAME}" "${out}/${NAME}"

# Mounts switchroot folder as boot folder if a hekate ID is given
if [[ -n ${HEKATE_ID} ]]; then
	mount --bind "${out}/switchroot/${DISTRO}" "${out}/${NAME}/boot/"

	if [ -e "${out}/switchroot/${DISTRO}/update.tar.gz" ]; then
		tar xhpf "${out}/switchroot/${DISTRO}/update.tar.gz" -C "${out}/${NAME}"
	fi

	if [ -e "${out}/switchroot/${DISTRO}/modules.tar.gz" ]; then
		tar xhpf "${out}/switchroot/${DISTRO}/modules.tar.gz" -C "${out}/${NAME}/lib/"
	fi
fi

# Add cache dir configuration
if [ -n "$CACHE_DIR" ]; then
	mkdir "${out}/cache" &> /dev/null || true
	mount --bind "${out}/cache" "${out}/${NAME}/${CACHE_DIR}" || exit
fi

# Copy build script
cp "$(dirname "${cwd}")/configs/${DEVICE}/files/${CHROOT_SCRIPT}" "${out}/${NAME}"

# Handle resolv.conf
cp --remove-destination --dereference /etc/resolv.conf "${out}/${NAME}/etc/resolv.conf"

# Actual chroot
arch-chroot "${out}/${NAME}" /bin/bash /"${CHROOT_SCRIPT}"

# Unmount switchroot boot dir
[[ -n ${HEKATE_ID} ]] && umount -l "${out}/${NAME}/boot"

# Unmount chroot dir
umount -l "${out}/${NAME}"

# Check lock status
if [[ ${lock} = 1 ]]; then
	# Get current lock count
	lock_count=$(sed 's/'${AARCH}' //g' "${out}/.lock")

	# If the current instance is the only one left for this binary, remove it
	if [[ ${lock_count} = 1 ]]; then
		# Remove lock on architecture
		sed -i '/'${AARCH}'*/d' "${out}/.lock"

		# Unregister binary if it wasn't set on script launch
		"${out}/downloadedFiles/register.sh" -- -r -p ${AARCH}
	else
		# Decrement lock count in lock file
		sed -i 's/'${AARCH}' '${lock_count}'/'${AARCH}' '$((lock_count-1))'/g' "${out}/.lock"
	fi
fi

# Remove lock file if empty, meaning no more instance is running.
[[ ! -s "${out}/.lock" ]] && rm -rf "${out}/.lock"

# Clean qemu emulation files
if [[ ${same_cpu_arch} = 0 ]]; then
	rm -rf "${out}/${NAME}/usr/bin/qemu-${AARCH}-static" "${out}/downloadedFiles/register.sh"
fi

rm -rf "${out}/${NAME}/${CHROOT_SCRIPT}" "${out}/cache"
