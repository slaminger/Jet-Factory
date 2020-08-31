#!/bin/bash
# DL_FILE.SH: Attempt to downlaod the image file
if [[ ! -e "${out}/downloadedFiles/${img%.*}" ]]; then
	if [[ -n "${URL}" ]]; then
		wget -q -nc --show-progress "${URL}" -P "${out}/downloadedFiles/"
	else
		echo "No URL found !";exit 1;
	fi
fi
