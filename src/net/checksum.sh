#!/usr/bin/bash

img_sig="${SIG##*/}"

# Download checksm if avalaible Check file integrity
wget -q --show-progress ${SIG} -P "$2/${img_sig}"

# MD5 Checksum
if [ ${SIG} =~ ".md5" ]; then
	md5sum --status -c "$2/${img_sig}"
else
	# SHA Checksum 
	# TODO
fi