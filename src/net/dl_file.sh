#!/usr/bin/bash
# DL_FILE.SH: Download a file and check it's integrity if possible

# Attempt to downlaod the image file
if [ ${URL} != "" ]; then
	wget -q --show-progress ${URL} -O "$2/${img}"
else
	echo "No URL found !"
	exit 1
fi