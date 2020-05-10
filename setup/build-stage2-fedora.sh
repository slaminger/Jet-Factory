#!/usr/bin/bash
uname -a

dnf -y update && dnf -y groupinstall 'Basic Desktop' 'Xfce Desktop' && dnf -y install `cat base-pkgs`
dnf -y remove xorg-x11-server-common iscsi-initiator-utils-iscsiuio iscsi-initiator-utils clevis-luks atmel-firmware kernel*
for pkg in `find /pkgs/*.rpm -type f`; do
	rpm -ivvh --force $pkg
done
dnf -y clean all && rm -r /pkgs

# TODO: Make kernel rpm
echo '\nexclude=linux-firmware kernel* xorg-x11-server-* xorg-x11-drv-ati xorg-x11-drv-armsoc xorg-x11-drv-nouveau xorg-x11-drv-ati xorg-x11-drv-qxl xorg-x11-drv-fbdev' >> /etc/dnf/dnf.conf

systemctl enable r2p bluetooth lightdm NetworkManager
sed -i 's/#keyboard=/keyboard=onboard/' /etc/lightdm/lightdm-gtk-greeter.conf

/usr/sbin/useradd -m fedora
/usr/sbin/usermod -aG video,audio,wheel fedora
echo "fedora:fedora" | /usr/bin/chpasswd && echo "root:root" | /usr/bin/chpasswd

/usr/sbin/ldconfig