#!/bin/bash

# Cut sig name from SIG URL
img_sig="${SIG##*/}"

# Download checksm if avalaible Check file integrity
wget -q --show-progress "${SIG}" -P "${out}/downloadedFiles/"

# MD5 Checksum
if [[ ${SIG} =~ .md5 ]]; then
	md5sum --status -c "${out}/downloadedFiles/${img_sig}"
	# SHA Checksum
	# TODO
fi
