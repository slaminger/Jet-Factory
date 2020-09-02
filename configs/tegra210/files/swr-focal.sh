#!/bin/bash

echo "Updating apt repos in rootfs"
sed -i 's/http:\/\/ports\.ubuntu\.com\/ubuntu-ports\//http:\/\/turul.canonical.com\//g' /etc/apt/sources.list
#echo 'deb https://repo.download.nvidia.com/jetson/common r32.4 main
#deb https://repo.download.nvidia.com/jetson/t210 r32.4 main' > /etc/apt/sources.list.d/nvidia-l4t-apt-source.list
echo TEGRA_CHIPID 0x21 > /etc/nv_boot_control.conf
mkdir -p /opt/nvidia/l4t-packages/
touch /opt/nvidia/l4t-packages/.nv-l4t-disable-boot-fw-update-in-preinstall
echo "Done!"

echo "Installing desktop packages"
export DEBIAN_FRONTEND=noninteractive
apt update
yes | unminimize
apt install -y openssh-server systemd wget gnupg nano sudo gnome-session gnome-terminal \
 ubuntu-desktop-minimal linux-firmware less bsdutils locales ||
(
 rm -rf /usr/share/dict/words.pre-dictionaries-common
 apt --fix-broken install
) # nicman23 says abracatabra ubuntu is shit
echo "Done!"

echo "Adding switchroot /nvidia key"
wget https://newrepo.switchroot.org/pubkey
apt-key add pubkey
rm pubkey
apt-key adv --fetch-key https://repo.download.nvidia.com/jetson/jetson-ota-public.asc
echo "Done!"


echo "Adding switchroot repo"
wget https://newrepo.switchroot.org/pool/unstable/s/switchroot-newrepo/switchroot-newrepo_1.1_all.deb
wget http://turul.canonical.com/pool/main/libf/libffi/libffi6_3.2.1-8_arm64.deb
dpkg -i switchroot-newrepo_1.1_all.deb libffi6_3.2.1-8_arm64.deb
rm switchroot-newrepo_1.1_all.deb libffi6_3.2.1-8_arm64.deb
echo 'force-overwrite' > /etc/dpkg/dpkg.cfg.d/sadface
apt update
apt dist-upgrade -y; apt install -y nintendo-switch-meta joycond
apt install -y nvidia-l4t-init nvidia-l4t-multimedia nvidia-l4t-oem-config \
 nvidia-l4t-3d-core nvidia-l4t-multimedia-utils nvidia-l4t-gstreamer \
 nvidia-l4t-firmware nvidia-l4t-xusb-firmware nvidia-l4t-configs \
 nvidia-l4t-tools nvidia-l4t-core nvidia-l4t-x11 nvidia-l4t-apt-source \
 nvidia-l4t-cuda nvidia-l4t-wayland
#rm /opt/nvidia/l4t-packages/.nv-l4t-disable-boot-fw-update-in-preinstall
echo "Done!"

echo "Making users"
useradd -m user
usermod -a -G sudo user
usermod -a -G video user
/usr/sbin/chpasswd << EOF
root:toor
user:user
EOF
echo "Done!"

echo "Fixing broken nvidia shit"
cat << EOF > /etc/systemd/system/upower.service
[Unit]
Description=Daemon for power management
Documentation=man:upowerd(8)

[Service]
Type=dbus
BusName=org.freedesktop.UPower
ExecStart=/usr/lib/upower/upowerd
Restart=on-failure

# Filesystem lockdown
ProtectSystem=strict
# Needed by keyboard backlight support
ProtectKernelTunables=false
ProtectControlGroups=true
ReadWritePaths=/var/lib/upower
StateDirectory=upower
ProtectHome=true
PrivateTmp=true

# Network
# PrivateNetwork=true would block udev's netlink socket
IPAddressDeny=any
RestrictAddressFamilies=AF_UNIX AF_NETLINK

# Execute Mappings
MemoryDenyWriteExecute=true

# Modules
ProtectKernelModules=true

# Real-time
RestrictRealtime=true

# Privilege escalation
NoNewPrivileges=true

# Capabilities
CapabilityBoundingSet=

# System call interfaces
LockPersonality=yes
SystemCallArchitectures=native
SystemCallFilter=@system-service
SystemCallFilter=ioprio_get

# Namespaces
PrivateUsers=no
RestrictNamespaces=no

# Locked memory
LimitMEMLOCK=0

[Install]
WantedBy=graphical.target
EOF

cat << EOF > /lib/firmware/brcm/brcmfmac4356-pcie.txt
NVRAMRev=662895
sromrev=11
boardrev=0x1250
boardtype=0x074a
boardflags=0x02400001
boardflags2=0xc0802000
boardflags3=0x00000108
macaddr=98:b6:e9:53:50:a3
ccode=0
regrev=0
antswitch=0
pdgain5g=4
pdgain2g=4
muxenab=0x10
wowl_gpio=0
wowl_gpiopol=0
swctrlmap_2g=0x11411141,0x42124212,0x10401040,0x00211212,0x000000ff
swctrlmap_5g=0x42124212,0x41114111,0x42124212,0x00211212,0x000000cf
swctrlmapext_2g=0x00000000,0x00000000,0x00000000,0x000000,0x003
swctrlmapext_5g=0x00000000,0x00000000,0x00000000,0x000000,0x003
phycal_tempdelta=50
papdtempcomp_tempdelta=20
fastpapdgainctrl=1
olpc_thresh=0
lowpowerrange2g=0
tworangetssi2g=1
lowpowerrange5g=0
tworangetssi5g=1
ed_thresh2g=-75
ed_thresh5g=-75
eu_edthresh2g=-75
eu_edthresh5g=-75
paprdis=0
femctrl=10
vendid=0x14e4
devid=0x43ec
manfid=0x2d0
nocrc=1
otpimagesize=502
xtalfreq=37400
rxchain=3
txchain=3
aa2g=3
aa5g=3
agbg0=2
agbg1=2
aga0=2
aga1=2
tssipos2g=1
extpagain2g=2
tssipos5g=1
extpagain5g=2
tempthresh=255
tempoffset=255
rawtempsense=0x1ff
pa2ga0=-181,5872,-700
pa2ga1=-180,6148,-728
pa2ga2=-193,3535,-495
pa2ga3=-201,3608,-499
pa5ga0=-189,5900,-717,-190,5874,-715,-189,5921,-718,-194,5812,-708
pa5ga1=-194,5925,-724,-196,5852,-718,-189,5858,-712,-196,5767,-707
pa5ga2=-187,3550,-504,-176,3713,-526,-189,3597,-505,-192,3532,-496
pa5ga3=-187,3567,-507,-187,3543,-506,-181,3589,-512,-187,3582,-508
subband5gver=0x4
pdoffsetcckma0=0x2
pdoffsetcckma1=0x2
pdoffset40ma0=0x3344
pdoffset80ma0=0x1133
pdoffset40ma1=0x3344
pdoffset80ma1=0x1133
maxp2ga0=76
maxp5ga0=74,74,74,74
maxp2ga1=76
maxp5ga1=74,74,74,74
cckbw202gpo=0x0000
cckbw20ul2gpo=0x0000
mcsbw202gpo=0x99644422
mcsbw402gpo=0x99644422
dot11agofdmhrbw202gpo=0x6666
ofdmlrbw202gpo=0x0022
mcsbw205glpo=0x88766663
mcsbw405glpo=0x88666663
mcsbw805glpo=0xbb666665
mcsbw205gmpo=0xd8666663
mcsbw405gmpo=0x88666663
mcsbw805gmpo=0xcc666665
mcsbw205ghpo=0xdc666663
mcsbw405ghpo=0xaa666663
mcsbw805ghpo=0xdd666665
mcslr5glpo=0x0000
mcslr5gmpo=0x0000
mcslr5ghpo=0x0000
sb20in40hrpo=0x0
sb20in80and160hr5glpo=0x0
sb40and80hr5glpo=0x0
sb20in80and160hr5gmpo=0x0
sb40and80hr5gmpo=0x0
sb20in80and160hr5ghpo=0x0
sb40and80hr5ghpo=0x0
sb20in40lrpo=0x0
sb20in80and160lr5glpo=0x0
sb40and80lr5glpo=0x0
sb20in80and160lr5gmpo=0x0
sb40and80lr5gmpo=0x0
sb20in80and160lr5ghpo=0x0
sb40and80lr5ghpo=0x0
dot11agduphrpo=0x0
dot11agduplrpo=0x0
temps_period=15
temps_hysteresis=15
rssicorrnorm_c0=4,4
rssicorrnorm_c1=4,4
rssicorrnorm5g_c0=1,1,3,1,1,2,1,1,2,1,1,2
rssicorrnorm5g_c1=3,3,4,3,3,4,3,3,4,2,2,3
initxidx2g=20
initxidx5g=20
btc_params84=0x8
btc_params95=0x0
btcdyn_flags=0x3
btcdyn_dflt_dsns_level=99
btcdyn_low_dsns_level=0
btcdyn_mid_dsns_level=22
btcdyn_high_dsns_level=24
btcdyn_default_btc_mode=5
btcdyn_btrssi_hyster=5
btcdyn_dsns_rows=1
btcdyn_dsns_row0=5,-120,0,-52,-72

EOF

systemctl enable upower
#apt clean

mkdir -p /usr/share/alsa/ucm/tegra-s/
ln -s /usr/share/alsa/ucm/tegra-snd-t210ref-mobile-rt565x/HiFi /usr/share/alsa/ucm/tegra-s/HiFi
sed 's/44100/48000/g' -i /etc/pulse/daemon.conf

echo "Done!"
