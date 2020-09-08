#!/bin/bash
echo "Installing DE, Nvidia drivers and configs..."
dnf config-manager --add-repo https://download.opensuse.org/repositories/home:/Azkali:/Switch-L4T/Fedora_32/home:Azkali:Switch-L4T.repo

dnf config-manager --add-repo https://download.opensuse.org/repositories/home:/Azkali:/branches:/home:/concyclic:/java/Fedora_32/home:Azkali:branches:home:concyclic:java.repo

dnf -y install @xfce-desktop-environment lightdm lightdm-gtk-greeter firefox onboard \
				langpacks-ja upower screen wpa_supplicant alsa-utils initial-setup-gui \
				alsa-ucm alsa-plugins-pulseaudio pulseaudio pulseaudio-module-x11 pulseaudio-utils \
				xorg-x11-xinit xorg-x11-drv-libinput xorg-x11-drv-wacom xorg-x11-drv-evdev \
				libnvmpi1_0_0 nvidia-l4t-init nvidia-l4t-multimedia nvidia-l4t-oem-config \
				nvidia-l4t-3d-core nvidia-l4t-multimedia-utils \
				nvidia-l4t-firmware nvidia-l4t-configs nvidia-l4t-tools \
				nvidia-l4t-core nvidia-l4t-x11 nvidia-l4t-cuda nvidia-l4t-wayland \
				https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm \
				https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm

dnf install --downloadonly -y switch-configs
rpm -ivvh --force /var/cache/dnf/home_Azkali_Switch-L4T-*/packages/switch-configs-*.aarch64.rpm

dnf -y clean all

# Userland configuration
systemctl enable bluetooth r2p NetworkManager upower lightdm initial-setup
sed -i 's/#keyboard=/keyboard=onboard/' /etc/lightdm/lightdm-gtk-greeter.conf
touch /.unconfigured
systemctl set-default graphical.target

# SDDM fix
# echo "SUBSYSTEM=="graphics", KERNEL=="fb[0-9]", TAG+="master-of-seat"" > /etc/udev/rules.d/69-nvidia-seat.rules

# Wi-Fi sleep issues fix
echo brcmfmac >etc/modules-load.d/brcmfmac.conf

# Audio Fix
sed 's/44100/48000/g' -i /etc/pulse/daemon.conf

# useradd -m -G video,audio,wheel -s /bin/bash fedora
# chown -R fedora:fedora /home/fedora
# echo "fedora:fedora" | chpasswd && echo "root:root" | chpasswd
