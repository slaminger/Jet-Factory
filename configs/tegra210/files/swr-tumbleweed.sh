#!/usr/bin/bash

echo "Installing switch drivers and configs"
zypper -n ar --refresh -p 90 https://download.opensuse.org/repositories/home:/Azkali:/Switch-L4T/openSUSE_Tumbleweed/home:Azkali:Switch-L4T.repo
zypper --gpg-auto-import-keys refresh
zypper -n in switch-meta
zypper -n clean -a
echo "Done!"

echo "Fixing boot stuff..."
rm -r /etc/X11/xorg.conf.d/20-fbdev.conf
echo "Done!"

echo "Configuring user..."
groupadd wheel
useradd -m -G wheel,video,audio,users -s /bin/bash suse
echo "suse:suse" | chpasswd && echo "root:root" | chpasswd
chown -R 1000:1000 /home/suse
echo "Done!"
