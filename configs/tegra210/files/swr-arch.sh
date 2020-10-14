#!/bin/bash

# Pre install configurations
## Workaround for flakiness of `pt` mirror.
sed -i 's/http:\/\/mirror.archlinuxarm.org/http:\/\/eu.mirror.archlinuxarm.org/g' /etc/pacman.d/mirrorlist

# Initialize pacman keyring
pacman-key --init
pacman-key --populate archlinuxarm

## Arch switchroot repository
echo -e "[switchrootarch]\nServer = https://archrepo.switchroot.org/" >> /etc/pacman.conf
curl https://newrepo.switchroot.org/pubkey > /tmp/pubkey
pacman-key --add /tmp/pubkey
pacman-key --lsign-key C9DDF6AA7BAC41CF6B893FB892813F6A23DB6DFC

pacman -Syy

# Installation
pacman -R linux-aarch64 --noconfirm

echo -e "\n\nBeginning packages installation!"
pacman -Syyu --noconfirm jetson-ffmpeg \
			tegra-bsp \
			onboard \
			wpa_supplicant \
			dialog \
			pulseaudio \
			pulseaudio-alsa \
			bluez \
			sudo \
			alsa-utils \
			dhcpcd \
			networkmanager \
			switch-configs \
			joycond-git

# Post install configurations
yes | pacman -Scc

## Audio fix
sed -i 's/.*default-sample-rate.*/default-sample-rate = 48000/' /etc/pulse/daemon.conf

## SDDM fix
echo "SUBSYSTEM=="graphics", KERNEL=="fb[0-9]", TAG+="master-of-seat"" > /etc/udev/rules.d/69-nvidia-seat.rules

systemctl enable r2p bluetooth NetworkManager
usermod -aG video,audio,wheel alarm
