#!/usr/bin/bash
# DL_FILE.SH: Download a file and check it's integrity if possible

# Attempt to downlaod the image file
if [ ${img} != "" ]; then
	wget -q --show-progress ${img} -O "$2/${img}"
else
	exit 1
fi