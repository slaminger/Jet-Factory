#!/usr/bin/bash

# Download checksm if avalaible Check file integrity
wget -q --show-progress ${img_sig} -O "$2/${img_sig}"

# MD5 Checksum
if [ ${img_sig} =~ ".md5" ]; then
	md5sum --status -c "$2/${img_sig}"
fi

# SHA Checksum 
# TODO