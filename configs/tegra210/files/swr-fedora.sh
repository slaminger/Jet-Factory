#!/bin/bash
echo "Installing DE, Nvidia drivers and configs..."
dnf config-manager --add-repo https://download.opensuse.org/repositories/home:/Azkali:/Switch-L4T/Fedora_32/home:Azkali:Switch-L4T.repo

dnf -y install @kde-desktop @deepin-desktop lightdm firefox onboard langpacks-ja upower screen wpa_supplicant alsa-utils \
				alsa-ucm alsa-plugins-pulseaudio pulseaudio pulseaudio-module-x11 pulseaudio-utils \
				xorg-x11-xinit xorg-x11-drv-libinput xorg-x11-drv-wacom xorg-x11-drv-evdev \
				nvidia-l4t-* libnvmpi1_0_0 \
				https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm \
				https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm
dnf --downloadonly -y switch-configs
rpm -ivvh --force switch-configs
dnf -y clean all

# Userland configuration
systemctl enable bluetooth lightdm r2p NetworkManager
sed -i 's/#keyboard=/keyboard=onboard/' /etc/lightdm/lightdm-gtk-greeter.conf
systemctl set-default graphical.target

# SDDM fix
echo "SUBSYSTEM=="graphics", KERNEL=="fb[0-9]", TAG+="master-of-seat"" > /etc/udev/rules.d/69-nvidia-seat.rules

# https://fluxcoil.net/hardwarerelated/nintendo_switch
echo brcmfmac >etc/modules-load.d/brcmfmac.conf

useradd -m -G video,audio,wheel -s /bin/bash fedora
chown -R fedora:fedora /home/fedora
echo "fedora:fedora" | chpasswd && echo "root:root" | chpasswd
