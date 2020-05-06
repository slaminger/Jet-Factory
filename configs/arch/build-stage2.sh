#!/bin/bash

# Pre install configurations
## Workaround for flakiness of `pt` mirror.
sed -i 's/mirror.archlinuxarm.org/eu.mirror.archlinuxarm.org/g' /etc/pacman.d/mirrorlist
echo -e "[switch]\nSigLevel = Optional\nServer = https://9net.org/l4t-arch/" >> /etc/pacman.conf
sed 's/.*default-sample-rate.*/default-sample-rate = 48000/' -i /etc/pulse/daemon.conf

# Configuring pacman
pacman-key --init
pacman-key --populate archlinuxarm

# Installation
## Removing linux-aarch64 as we won't be needing this
pacman -R linux-aarch64 --noconfirm

i=5
echo -e "\n\nBeginning packages installation!\nRetry attempts left: ${i}"
until [[ ${i} == 0 ]] || pacman -Syu `cat /base-pkgs` --noconfirm; do
	pacman -Syu `cat /base-pkgs` --noconfirm
	echo -e "\n\nPackages installation failed, retrying!\nRetry attempts left: ${i}"
	let --i
done


for pkg in `find /pkgs/*.pkg.* -type f`; do
	pacman -U $pkg --noconfirm
done

yes | pacman -Scc

# Post install configurations
systemctl enable r2p bluetooth lightdm NetworkManager

echo brcmfmac > /etc/suspend-modules.conf
sed -i 's/#keyboard=/keyboard=onboard/' /etc/lightdm/lightdm-gtk-greeter.conf

usermod -aG video,audio,wheel alarm

ldconfig

echo "Exit chroot."
