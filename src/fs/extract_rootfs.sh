#!/usr/bin/bash
# EXTRACT_ROOTFS.SH: Extract Rootfs using libguetsfs-tools
# Handles .raw disk image
if [ ${img} =~ ".raw.xz" ]; then
	# Uncompress (xz format) the filesystem archive file
	[[ $(file -b --mime-type "${img}") == "application/x-xz" ]] && unxz "${img}"

	# https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html
	extracted_img=${img%.*}

	# Search for a partition labeled root (lvm2 handling) and extract it
	partition=$(virt-filesystems -a ${extracted_img} | grep "root")

	# Extract the partition from the image file
	guestfish --ro -a ${extracted_img} -m ${partition} tgz-out / tmp.tar.gz

	# Reassign img to previously created tar archive
	img=tmp.tar.gz
fi

# Handles tar archive
if [ ${img} =~ ".tar" ]; then
	tar xf ${img} $1
else
	echo "Unrecognzied format, exiting !"
	exit 1
fi