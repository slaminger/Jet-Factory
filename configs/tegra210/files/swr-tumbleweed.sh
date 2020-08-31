#!/usr/bin/bash

echo "Installing XFCE, Nvidia drivers and switch config..."
zypper -n ar --refresh -p 90 https://download.opensuse.org/repositories/home:/Azkali:/Switch-L4T/openSUSE_Tumbleweed/home:Azkali:Switch-L4T.repo
zypper --gpg-auto-import-keys refresh
zypper -n in libnvmpi1_0_0 nvidia-l4t-* 
# zypper -n in --oldpackage xorg-x11-server-1.19.6 xorg-x11-server-extra-1.19.6 xorg-x11-server-source-1.19.6
zypper -n in switch-configs
zypper -n clean all
systemctl enable r2p bluetooth NetworkManager
echo "Done!"

echo "Fixing boot stuff..."
rm -r /etc/X11/xorg.conf.d/20-fbdev.conf 
update-alternatives --set default-displaymanager /usr/lib/X11/displaymanagers/lightdm
echo "keyboard=onboard" >> /etc/lightdm/lightdm-gtk-greeter.conf
echo "Done!"

echo "Configuring user..."
groupadd wheel
useradd -m -G wheel,video,audio,users -s /bin/bash suse
echo "suse:suse" | chpasswd && echo "root:root" | chpasswd
chown -R 1000:1000 /home/suse
sed -i 's/#%wheel        ALL=(ALL) ALL/%wheel        ALL=(ALL) ALL/g' /etc/sudoers
echo "Done!"
