#!/bin/bash
# MAKEIMG.SH : Create rootfs image file

# Set Image name
guestfs_img="switchroot-${DISTRO}.img"

# Clean previously made image file
[[ -f "${guestfs_img}" ]] && rm -rf "${guestfs_img}"

# Convert to hekate format or create image
if [[ "${HEKATE}" == "true" ]]; then
	echo "Creating hekate installable partition..."
	source "${cwd}/fs/hekate.sh"
else
	# Create image
	virt-make-fs --type=ext4 --format=raw --size=+512MB "${out}/${NAME}/" ${guestfs_img}
	zerofree -n ${guestfs_img}
fi

# Clean unneeded files
rm -rf "${out}/${NAME}/"
