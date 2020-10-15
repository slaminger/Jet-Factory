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

	rootfs="$(virt-filesystems -a "${extracted_img}" | grep "root")"
	if [[ -z ${rootfs} ]]; then
		# Search for a partition labeled root (lvm2 handling) and extract it
		parts_sizes="$(virt-filesystems -la "${extracted_img}" | awk 'NR > 1 {print $5}')"
		i=0

		for p_size in $parts_sizes; do
			if [[ -z $biggest ]] || [[ ${!biggest} -lt ${p_size} ]]; then
				declare size_$i=$p_size
				biggest=size_$i
				((i=i+1))
			fi
		done

		# Get biggest partition 
		rootfs=${!biggest}
	fi

	# Store tar archive as img
	img="${out}/downloadedFiles/tmp.tar.gz"

	# Extract the partition from the image file
	guestfish --ro -a "${extracted_img}" -m "${rootfs}" tgz-out / "${img}"

	# Extract tmp.tar.gz
	tar xpf "${img}" -C "${out}/${NAME}"

	# Cleanup tmp.tar.gz
	rm -rf "${out}/downloadedFiles/tmp.tar.gz"
elif [[ "${img}" =~ .iso ]]; then
	echo "Iso not implemented yet.."
	exit 1
elif [[ "${img}" =~ .tbz2 ]]; then
	tar xpf "${img}" -C "${out}/${NAME}"
elif [[ "${img}" =~ .tar ]]; then
	bsdtar xpf "${img}" -C "${out}/${NAME}"
else
	echo "Unrecognzied format, exiting !"
	exit 1
fi
