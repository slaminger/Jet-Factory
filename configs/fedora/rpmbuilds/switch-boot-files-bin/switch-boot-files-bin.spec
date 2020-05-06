# Maintainer: Ezekiel Bethel <zek@9net.org>
Name:			switch-boot-files-bin
Version:		R32
Release:		4.2
BuildArch:		noarch
Source0:		l4t-fedora.ini
Source1:		mkscr
Source2:		boot.scr.txt
Source3:		coreboot.rom
License:        GPLv3+
Summary:		Switch boot files
BuildRequires:	uboot-tools

%description
	Switch boot files

%prep
	mkdir -p %buildroot/boot/bootloader/ini %buildroot/boot/switchroot/fedora/

%install
	install %SOURCE0 %buildroot/boot/bootloader/ini/switchroot-fedora.ini
	install %SOURCE1 %SOURCE2 %SOURCE3 %buildroot/boot/switchroot/fedora/
	mkimage -A arm -T script -O linux -d %SOURCE2 %buildroot/boot/switchroot/fedora/boot.scr

	mkdir -p %buildroot/boot/l4t-fedora/
	cp -r bootloader %buildroot/boot/
	cp -r l4t-fedora/{boot.scr,coreboot.rom} %buildroot/boot/l4t-fedora

%files
/boot/*

%clean
rm -rf %buildroot