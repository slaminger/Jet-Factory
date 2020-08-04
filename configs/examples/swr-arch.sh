#!/bin/bash

# Pre install configurations
## Workaround for flakiness of `pt` mirror.
sed -i 's/http://http://mirror.archlinuxarm.org/http://eu.mirror.archlinuxarm.org/g' /etc/pacman.d/mirrorlist

# Initialize pacman keyring
pacman-key --init
pacman-key --populate archlinuxarm

## Arch switchroot repository
echo -e "[switchrootarch]\nServer = https://archrepo.switchroot.org/" >> /etc/pacman.conf
curl https://newrepo.switchroot.org/pubkey > /tmp/pubkey
pacman-key --add /tmp/pubkey
pacman-key --lsign-key C9DDF6AA7BAC41CF6B893FB892813F6A23DB6DFC

## Audio fix
# sed 's/.*default-sample-rate.*/default-sample-rate = 48000/' -i /etc/pulse/daemon.conf

# Installation
pacman -R linux-aarch64 --noconfirm

echo -e "\n\nBeginning packages installation!"
# pacman -Syyu joycond-git switch-boot-files-bin systemd-suspend-modules xorg-server-tegra switch-config tegra-bsp linux-tegra --noconfirm
pacman -Syyu --noconfirm jetson-ffmpeg tegra-ffmpeg tegra-bsp \
			xorg-xrandr xorg-xinput xorg-xinit onboard \
			wpa_supplicant dialog pulseaudio pulseaudio-alsa \
			bluez sudo lightdm lightdm-gtk-greeter plasma \
			kde-applications plasma-wayland-session alsa-utils \
			netctl dhcpcd networkmanager

# Post install configurations
yes | pacman -Scc
systemctl enable r2p bluetooth lightdm NetworkManager
echo brcmfmac > /etc/suspend-modules.conf
sed -i 's/#keyboard=/keyboard=onboard/' /etc/lightdm/lightdm-gtk-greeter.conf
usermod -aG video,audio,wheel alarm