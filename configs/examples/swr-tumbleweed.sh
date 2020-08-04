#!/usr/bin/bash

echo "Installing XFCE, Nvidia drivers and switch config..."
zypper up -y && zypper -t pattern xfce && zypper -n clean all
systemctl enable r2p bluetooth NetworkManager
echo "Done!"

echo "Fixing boot stuff..."
rm -r /etc/X11/xorg.conf.d/20-fbdev.conf 
update-alternatives --set default-displaymanager /usr/lib/X11/displaymanagers/lightdm
echo "keyboard=onboard" >> /etc/lightdm/lightdm-gtk-greeter.conf
echo "Done!"

echo "Configuring user..."
useradd -m -G wheel,video,audio,users -s /bin/bash suse
echo "suse:suse" | chpasswd && echo "root:root" | chpasswd
sed -i 's/#%wheel        ALL=(ALL) ALL/%wheel        ALL=(ALL) ALL/g' /etc/sudoers
echo "Done!"