#!/bin/bash
echo "Updating apt repos in rootfs"
sed -i 's/http:\/\/ports\.ubuntu\.com\/ubuntu-ports\//http:\/\/turul.canonical.com\//g' /etc/apt/sources.list
# echo 'deb https://repo.download.nvidia.com/jetson/common r32.4 main
# deb https://repo.download.nvidia.com/jetson/t210 r32.4 main' > /etc/apt/sources.list.d/nvidia-l4t-apt-source.list
echo TEGRA_CHIPID 0x21 > /etc/nv_boot_control.conf
mkdir -p /opt/nvidia/l4t-packages/
touch /opt/nvidia/l4t-packages/.nv-l4t-disable-boot-fw-update-in-preinstall
echo "Done!"

echo "Installing desktop packages"
export DEBIAN_FRONTEND=noninteractive
apt update
yes | unminimize
echo "Done!"

echo "Adding switchroot repo"
wget https://newrepo.switchroot.org/pubkey
apt-key add pubkey
rm pubkey
apt-key adv --fetch-key https://repo.download.nvidia.com/jetson/jetson-ota-public.asc
wget https://newrepo.switchroot.org/pool/unstable/s/switchroot-newrepo/switchroot-newrepo_1.1_all.deb
dpkg -i switchroot-newrepo_1.1_all.deb
rm switchroot-newrepo_1.1_all.deb
echo "Done"

echo "Installing Tegra210 BSP, Switch config and Joycond"
apt update -y && apt install -y nintendo-switch-meta joycond \
	nvidia-l4t-init nvidia-l4t-multimedia nvidia-l4t-oem-config \
	nvidia-l4t-3d-core nvidia-l4t-multimedia-utils nvidia-l4t-gstreamer \
	nvidia-l4t-firmware nvidia-l4t-xusb-firmware nvidia-l4t-configs \
	nvidia-l4t-tools nvidia-l4t-core nvidia-l4t-x11 nvidia-l4t-apt-source \
	nvidia-l4t-cuda nvidia-l4t-wayland nvidia-l4t-core xxd || true
	
apt clean
sed 's/44100/48000/g' -i /etc/pulse/daemon.conf
echo "Done"
# rm /opt/nvidia/l4t-packages/.nv-l4t-disable-boot-fw-update-in-preinstall

# mkdir -p /usr/share/alsa/ucm/tegra-s/
# ln -s /usr/share/alsa/ucm/tegra-snd-t210ref-mobile-rt565x/HiFi /usr/share/alsa/ucm/tegra-s/HiFi
