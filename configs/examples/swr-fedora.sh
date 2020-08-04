#!/usr/bin/bash

echo "Installing XFCE, Nvidia drivers and switch config..."
dnf -y update && dnf -y install @kde-desktop lightdm git onboard langpacks-ja upower screen \
 								wpa_supplicant alsa-utils alsa-ucm alsa-plugins-pulseaudio pulseaudio pulseaudio-module-x11 \
  								pulseaudio-utils xorg-x11-xinit xorg-x11-drv-libinput xorg-x11-drv-wacom xorg-x11-drv-evdev \
								https://kojipkgs.fedoraproject.org//vol/fedora_koji_archive02/packages/xorg-x11-server/1.19.6/7.fc28/aarch64/xorg-x11-server-Xdmx-1.19.6-7.fc28.aarch64.rpm \
								https://kojipkgs.fedoraproject.org//vol/fedora_koji_archive02/packages/xorg-x11-server/1.19.6/7.fc28/aarch64/xorg-x11-server-Xephyr-1.19.6-7.fc28.aarch64.rpm \
								https://kojipkgs.fedoraproject.org//vol/fedora_koji_archive02/packages/xorg-x11-server/1.19.6/7.fc28/aarch64/xorg-x11-server-Xnest-1.19.6-7.fc28.aarch64.rpm \
								https://kojipkgs.fedoraproject.org//vol/fedora_koji_archive02/packages/xorg-x11-server/1.19.6/7.fc28/aarch64/xorg-x11-server-Xorg-1.19.6-7.fc28.aarch64.rpm \
								https://kojipkgs.fedoraproject.org//vol/fedora_koji_archive02/packages/xorg-x11-server/1.19.6/7.fc28/aarch64/xorg-x11-server-Xvfb-1.19.6-7.fc28.aarch64.rpm \
								https://kojipkgs.fedoraproject.org//vol/fedora_koji_archive02/packages/xorg-x11-server/1.19.6/7.fc28/aarch64/xorg-x11-server-Xwayland-1.19.6-7.fc28.aarch64.rpm \
								https://kojipkgs.fedoraproject.org//vol/fedora_koji_archive02/packages/xorg-x11-server/1.19.6/7.fc28/aarch64/xorg-x11-server-common-1.19.6-7.fc28.aarch64.rpm \
								https://kojipkgs.fedoraproject.org//vol/fedora_koji_archive02/packages/xorg-x11-server/1.19.6/7.fc28/aarch64/xorg-x11-server-devel-1.19.6-7.fc28.aarch64.rpm \
								https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm \
								https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm

dnf -y remove xorg-x11-server-common iscsi-initiator-utils-iscsiuio iscsi-initiator-utils clevis-luks atmel-firmware kernel*
dnf -y clean all
echo "Done!"

# TODO: Make kernel rpm
echo '\nexclude=linux-firmware kernel* xorg-x11-server-* xorg-x11-drv-ati xorg-x11-drv-armsoc xorg-x11-drv-nouveau xorg-x11-drv-ati xorg-x11-drv-qxl xorg-x11-drv-fbdev' >> /etc/dnf/dnf.conf
systemctl enable bluetooth lightdm r2p NetworkManager
sed -i 's/#keyboard=/keyboard=onboard/' /etc/lightdm/lightdm-gtk-greeter.conf
systemctl set-default graphical.target
useradd -m -G video,audio,wheel -s /bin/bash fedora
echo "fedora:fedora" | chpasswd && echo "root:root" | chpasswd