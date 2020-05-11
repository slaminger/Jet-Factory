#!/bin/bash
# Source : https://ilhicas.com/2018/08/08/bash-script-to-install-packages-multiple-os.html
declare -A osInfo;
osInfo[/etc/debian_version]="apt-get update -y && apt-get install -y"
osInfo[/etc/alpine-release]="apk"
osInfo[/etc/centos-release]="yum update && yum install"
osInfo[/etc/fedora-release]="dnf update && dnf install"
osInfo[/etc/arch-release]="pacman -Syu"

for f in ${!osInfo[@]}
do
    if [[ -f $f ]];then
        package_manager=${osInfo[$f]}
    fi
done

echo $package_manager