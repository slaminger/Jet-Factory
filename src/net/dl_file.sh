#!/bin/bash
# DL_FILE.SH: Attempt to downlaod the image file
if [[ -n "${URL}" ]];then wget -q -nc --show-progress "${URL}" -P "${out}/downloadedFiles/";
	else echo "No URL found !";exit 1; fi
