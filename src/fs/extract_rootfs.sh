#!/bin/bash
# EXTRACT_ROOTFS.SH: Extract Rootfs using libguetsfs-tools
img="${out}/downloadedFiles/${img}"

# Handles .raw disk image
if [[ "${img}" == *.raw.xz ]]; then
	# https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html
	extracted_img="${img%.*}"

	# Uncompress (xz format) the filesystem archive file
	[[ "$(file -b --mime-type "${img}")" == "application/x-xz" ]] && \
		[[ ! -e "${extracted_img}" ]] && unxz "${img}"

	# Search for a partition labeled root (lvm2 handling) and extract it
	partition="$(virt-filesystems -a "${extracted_img}" | grep "root")"

	# Store tar archive as img
	img="${out}/downloadedFiles/tmp.tar.gz"

	# Extract the partition from the image file
	guestfish --ro -a "${extracted_img}" -m "${partition}" tgz-out / "${img}"

	# Cleanup temporary file
elif [[ "${img}" == *.iso ]]; then
	echo ""
fi

# Handles tar archive
if [[ "${img}" =~ .tar ]]; then
	tar xf "${img}" -C "${out}/${NAME}"
else
	echo "Unrecognzied format, exiting !"
	exit 1
fi

[[ -e "${out}/downloadedFiles/tmp.tar.gz" ]] && rm "${out}/downloadedFiles/tmp.tar.gz"
