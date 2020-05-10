#!/usr/bin/bash
uname -a

zypper -y update && zypper -y install `cat base-pkgs`
# for pkg in `find /pkgs/*.rpm -type f`; do
# 	zypper -ivvh --force $pkg
# done
# rm -r /pkgs
zypper -y clean all 

systemctl enable r2p bluetooth NetworkManager
sed -i 's/#keyboard=/keyboard=onboard/' /etc/lightdm/lightdm-gtk-greeter.conf

/usr/sbin/ldconfig
