#!/bin/bash
echo "Installing XFCE, Nvidia drivers and configs..."
dnf -y remove iscsi-initiator-utils-iscsiuio iscsi-initiator-utils clevis-luks atmel-firmware kernel*
dnf config-manager --add-repo https://download.opensuse.org/repositories/home:/Azkali:/Switch-L4T/Fedora_32/home:Azkali:Switch-L4T.repo

dnf -y install @xfce-desktop lightdm onboard langpacks-ja upower screen wpa_supplicant alsa-utils \
				alsa-ucm alsa-plugins-pulseaudio pulseaudio pulseaudio-module-x11 pulseaudio-utils \
				xorg-x11-xinit xorg-x11-drv-libinput xorg-x11-drv-wacom xorg-x11-drv-evdev \
				https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm \
				https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm

# Remove currently installed kernel and xorg-x11-server 1.20
rpm -e --noscripts --nodeps xorg-x11-server-Xorg xorg-x11-server-common kernel-core-5.7.14 kernel-core-5.6.6 kernel-headers kernel-modules kernel-modules-extra
dnf install -y xorg-x11-server*-1.19.6 nvidia-l4t-* libnvmpi1_0_0
dnf -y clean all

# Exclude xorg and kernel
echo -e '\nexclude=linux-firmware kernel* xorg-x11-server-common xorg-x11-server-Xorg xorg-x11-drv-ati xorg-x11-drv-armsoc xorg-x11-drv-nouveau xorg-x11-drv-ati xorg-x11-drv-qxl xorg-x11-drv-fbdev' >> /etc/dnf/dnf.conf

# Userland configuration
systemctl enable bluetooth lightdm r2p NetworkManager
sed -i 's/#keyboard=/keyboard=onboard/' /etc/lightdm/lightdm-gtk-greeter.conf
systemctl set-default graphical.target

# https://fluxcoil.net/hardwarerelated/nintendo_switch
echo brcmfmac >etc/modules-load.d/brcmfmac.conf

for i in auditd smartd pcscd ModemManager multipathd mdmonitor \
         dmraid-activation initial-setup lvm2-monitor zram-swap \
         plymouth-start lm_sensors udisks2 ; do
    systemctl disable $i
done

### tuning, seen in the L4T scripts
cat >/etc/rc.local<<EOT
#!/usr/bin/bash
echo 2048 > /sys/block/mmcblk0/queue/read_ahead_kb
echo 0 > "/proc/sys/vm/lazy_vfree_pages"
EOT
chmod +x /etc/rc.local
/etc/rc.local
ln -s /etc/rc.local /etc/rc.d/rc.local

useradd -m -G video,audio,wheel -s /bin/bash fedora
chown -R 1000:1000 /home/fedora
echo "fedora:fedora" | chpasswd && echo "root:root" | chpasswd

echo "Done!"
