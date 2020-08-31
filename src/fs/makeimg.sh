#!/bin/bash
# MAKEIMG.SH : Create rootfs image file

# Set Image name
guestfs_img="switchroot-${DISTRO}.img"

# Clean previously made image file
[[ -f "${guestfs_img}" ]] && rm "${guestfs_img}"

# Get build directory size
dir_size="$(du -hs -BM "${out}/${NAME}/" | awk '{print int($1+512);}')"

# Set extra space to 0 by default
aligned_size_extra=0

# If the build directory size modulo 4 isn't 0 then necessary extra space for alignement
[[ $((${dir_size%?}%4)) != 0 ]] && aligned_size_extra=$(awk '{print (4 - int($1%4));}' <<< ${dir_size})

# Convert to hekate format or create image
if [[ "${HEKATE}" == "true" ]]; then
	echo "Creating hekate installable partition..."
	source "${cwd}/fs/hekate.sh"
else
	# Create 4MB aligned image
	virt-make-fs --type=ext4 --format=raw --size=+$((512+${aligned_size_extra}))M "${out}/${NAME}/" ${guestfs_img}
fi

# Clean unneeded files
rm -r "${out}/${NAME}/"
