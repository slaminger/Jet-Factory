#!/usr/bin/bash
# DL_FILE.SH: Download a file and check it's integrity if possible

# Attempt to downlaod the image file
if [[ ${URL} != "" ]]; then
	wget -q -nc --show-progress ${URL} -P "${out}/downloadedFiles/"
else
	echo "No URL found !"
	exit 1
fi