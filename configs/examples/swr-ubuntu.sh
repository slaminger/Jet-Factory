#!/bin/bash

echo "Updating apt repos in rootfs"
sed -i 's/http:\/\/ports\.ubuntu\.com\/ubuntu-ports\//http:\/\/turul.canonical.com\//g' /etc/apt/sources.list
echo "Done!"


echo "Installing desktop packages"
yes | unminimize
echo "Done!"

echo "Adding switchroot key"
wget https://newrepo.switchroot.org/pubkey
apt-key add pubkey
rm pubkey
echo "Done!"

echo "Adding switchroot repo"
export DEBIAN_FRONTEND=noninteractive
wget https://newrepo.switchroot.org/pool/unstable/s/switchroot-newrepo/switchroot-newrepo_1.1_all.deb
dpkg -i switchroot-newrepo_1.1_all.deb
rm switchroot-newrepo_1.1_all.deb
echo "Done!"

echo "Backing up default bluetooth config"
cp /lib/systemd/system/bluetooth.service .
echo "Done!"

echo "Restoring up default bluetooth config"
mv ./bluetooth.service /lib/systemd/system/ -f
echo "Done!"

git clone https://github.com/Azkali/L4T-Packages-Repository/
SCRIPT_DIR="$PWD/L4T-Packages-Repository/"

echo "Patching USB Gadget script"
sed -i 's+cp -r /proc/device-tree/chosen/plugin-manager \"${mntpoint}/version/plugin-manager\"+\#cp -r /proc/device-tree/chosen/plugin-manager \"${mntpoint}/version/plugin-manager\"+g' opt/nvidia/l4t-usb-device-mode/nv-l4t-usb-device-mode-start.sh
echo "Done!"

echo "Patching ubiquity"
mkdir -p /usr/lib/ubiquity/dm-scripts/oem
cp $SCRIPT_DIR/files/ubiquity-oem-script /usr/lib/ubiquity/dm-scripts/oem/switch-randr
patch /usr/bin/ubiquity-dm $SCRIPT_DIR/files/ubiquity-force-enable-onboard.patch
patch /usr/bin/ubiquity-dm $SCRIPT_DIR/files/ubiquity-use-openbox.patch
patch /usr/sbin/oem-config-remove-gtk $SCRIPT_DIR/files/oem-config-uninstall-openbox.patch
echo "Done!"

echo "Applying switchroot customizations"
patch /etc/xdg/autostart/nvbackground.sh $SCRIPT_DIR/files/nvbackground.patch
cp $SCRIPT_DIR/files/Switchroot_Wallpaper.png /usr/share/backgrounds/
rm -f /etc/skel/Desktop/*
echo "Done"

echo "Installing display script"
cp $SCRIPT_DIR/files/nintendo-switch-display.desktop /usr/share/gdm/greeter/autostart/
cp $SCRIPT_DIR/files/nintendo-switch-display.desktop /etc/xdg/autostart/
echo "Done!"