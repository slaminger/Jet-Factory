#!/bin/bash
echo "Installing switch drivers and configs..."
dnf config-manager --add-repo https://download.opensuse.org/repositories/home:/Azkali:/Switch-L4T/Fedora_32/home:Azkali:Switch-L4T.repo

dnf -y install switch-meta \
	https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm \
	https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm

dnf -y clean all
rpm -e --noscripts kernel-core-5.6* kernel-core-5.8* kernel-modules* linux-firmware kernel-headers*

# Userland configuration
useradd -m -G video,audio,wheel -s /bin/bash fedora
chown -R fedora:fedora /home/fedora
echo "fedora:fedora" | chpasswd && echo "root:root" | chpasswd
