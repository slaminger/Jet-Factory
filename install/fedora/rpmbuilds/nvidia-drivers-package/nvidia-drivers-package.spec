# Maintainer: Your Name <youremail@domain.com>
%global nv_version "r32_Release_v4.2"
%global nv_dir   %{nv_name}-%{version}

Name:		nvidia-l4t
Version:	R32
Release:	4.2
License:	GPL
BuildArch:	noarch
Summary:	Nvidia L4T drivers
URL:		https://developer.nvidia.com/embedded/L4T/%{nv_version}/t210ref_release_aarch64/Tegra210_Linux_%{version}.%{release}_aarch64.tbz2
Requires:	cairo gstreamer1 pango libEGL libXext libX11 gstreamer1-plugins-base


%description
  	Nvidia L4T drivers for Jetson nano

%package 3d-core
	Summary: Nvidia L4T 3d-core drivers
	Requires: >= %{version}-%{release}
  
%description 3d-core
	3D Core drivers for Jetson nano

%package core
	Summary: Nvidia L4T core drivers
	Requires: >= %{version}-%{release}

%description core
	Core drivers for Jetson nano

%package configs
	Summary: Nvidia L4T configs
	Requires: >= %{version}-%{release}

%description configs
	Configs for Jetson nano

%package cuda
	Summary: Nvidia L4T cuda drivers
	Requires: >= %{version}-%{release}

%description cuda
	Cuda drivers for Jetson nano


%package firmware
	Summary: Nvidia L4T firmware
	Requires: >= %{version}-%{release}

%description firmware
	Firmware for Jetson nano

%package gstreamer
	Summary: Nvidia L4T gstreamer drivers
	Requires: >= %{version}-%{release}

%description gstreamer
	gstreamer drivers for Jetson nano

%package init
	Summary: Nvidia L4T init
	Requires: >= %{version}-%{release}

%description init
	init for Jetson nano

%package initrd
	Summary: Nvidia L4T initrd
	Requires: >= %{version}-%{release}

%description initrd
	initrd for Jetson nano

%package jetson-io
	Summary: Nvidia L4T jetson-io drivers
	Requires: >= %{version}-%{release}

%description jetson-io
	jetson-io drivers for Jetson nano

%package multimedia
	Summary: Nvidia L4T multimedia drivers
	Requires: >= %{version}-%{release}

%description multimedia
	multimedia drivers for Jetson nano

%package multimedia-utils
	Summary: Nvidia L4T multimedia-utils drivers
	Requires: >= %{version}-%{release}

%description multimedia-utils
	multimedia-utils drivers for Jetson nano

%package oem-config
	Summary: Nvidia L4T oem-config
	Requires: >= %{version}-%{release}

%description oem-config
	oem-config for Jetson nano

%package tools
	Summary: Nvidia L4T tools
	Requires: >= %{version}-%{release}

%description tools
	tools for Jetson nano

%package wayland
	Summary: Nvidia L4T wayland drivers
	Requires: >= %{version}-%{release}

%description wayland
	wayland drivers for Jetson nano

%package weston
	Summary: Nvidia L4T weston drivers
	Requires: >= %{version}-%{release}

%description weston
	weston drivers for Jetson nano

%package x11
	Summary: Nvidia L4T x11 drivers
	Requires: >= %{version}-%{release}

%description x11
	x11 drivers for Jetson nano

%package xusb-firmware
	Summary: Nvidia L4T xusb-firmware
	Requires: >= %{version}-%{release}

%description xusb-firmware
	xusb-firmware for Jetson nano

%prep
	rm -rf %{nv_dir}
	wget %{url} -P %{nv_dir}
	cd %{nv_dir}
	tar xvf Tegra210_Linux_%{version}.%{release}_aarch64.tbz2
	rm Tegra210_Linux_%{version}.%{release}_aarch64.tbz2

%build
	cd %{nv_dir}
	tar xvf Linux_for_Tegra/nv_tegra/nvidia_drivers.tbz2
	tar xvf Linux_for_Tegra/nv_tegra/config.tbz2

%post
	/sbin/ldconfig

%files 3d-core


%files core

%files configs

%files cuda

%file init

%file initrd

%files jetson-io

%file multimedia

%file multimedia-utils

%file oem-config

%file tools

%file wayland

%file weston

%file x11

%files xusb-firmware


/usr/*
/etc/*
/opt/*
/var/*

%clean
rm -rf %{buildroot}