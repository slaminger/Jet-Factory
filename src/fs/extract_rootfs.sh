#!/bin/bash
# EXTRACT_ROOTFS.SH: Extract Rootfs using libguetsfs-tools
img="${out}/downloadedFiles/${img}"

# Handles .raw disk image
if [[ ${img} =~ ".raw.xz" ]]; then
	# Uncompress (xz format) the filesystem archive file
	[[ $(file -b --mime-type "${img}") == "application/x-xz" ]] && unxz "${img}"

	# https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html
	extracted_img=${img%.*}

	# Search for a partition labeled root (lvm2 handling) and extract it
	partition=$(virt-filesystems -a ${extracted_img} | grep "root")

	# Store tar archive as img 
	img="${out}/downloadedFiles/tmp.tar.gz"

	# Extract the partition from the image file
	guestfish --ro -a ${extracted_img} -m ${partition} tgz-out / ${img}
fi

# Handles tar archive
if [[ ${img} =~ ".tar" ]]; then
	tar xf ${img} -C "${out}/${NAME}"
else
	echo "Unrecognzied format, exiting !"
	exit 1
fi
