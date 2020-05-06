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
Requires:   libXext libX11 libc6 libgcc1 Mesa-libEGL1 libstdc++6 Mesa-libGLESv2 libexpat1 libglvnd0 libasound2 libwayland-server libwayland-client libxkbcommon0
Requires:   libXdmcp libxcb1 gstreamer1-plugins-bad-free gstreamer1-plugins-base gstreamer1 libcairo2 glib2-devel pcre-devel libwayland-cursor libwayland-egl
%description
  	Nvidia L4T drivers for Jetson nano

%package 3d-core
	Summary: Nvidia L4T 3d-core drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: nvidia-l4t-firmware >= %{version}-%{release}
	Requires: nvidia-l4t-init >= %{version}-%{release}
	Requires: nvidia-l4t-wayland >= %{version}-%{release}
	Requires: nvidia-l4t-x11 >= %{version}-%{release}
	Requires: libXext6

	%description 3d-core
		3D Core drivers for Jetson nano

%package camera
	Summary: Nvidia L4T core drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: nvidia-l4t-3d-core >= %{version}-%{release}
	Requires: nvidia-l4t-cuda >= %{version}-%{release}
	Requires: nvidia-l4t-multimedia >= %{version}-%{release}
	Requires: nvidia-l4t-multimedia-utils >= %{version}-%{release}

	%description camera
		Core drivers for Jetson nano

%package configs
	Summary: Nvidia L4T configs
	Requires: nvidia-l4t-core >= %{version}-%{release}

	%description configs
		Configs for Jetson nano

%package core
	Summary: Nvidia L4T core drivers
	Requires: >= %{version}-%{release}

	%description core
		Core drivers for Jetson nano

%package cuda
	Summary: Nvidia L4T cuda drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: nvidia-l4t-3d-core >= %{version}-%{release}

	%description cuda
		Cuda drivers for Jetson nano


%package firmware
	Summary: Nvidia L4T firmware
	Requires: nvidia-l4t-core >= %{version}-%{release}

	%description firmware
		Firmware for Jetson nano

%package graphics-demos
	Summary: Nvidia L4T graphics demos
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: nvidia-l4t-3d-core >= %{version}-%{release}
	Requires: nvidia-l4t-weston >= %{version}-%{release}

	%description graphics-demos
		graphics-demos for Jetson nano

%package gstreamer
	Summary: Nvidia L4T gstreamer drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: nvidia-l4t-camera >= %{version}-%{release}
	Requires: nvidia-l4t-cuda >= %{version}-%{release}
	Requires: nvidia-l4t-multimedia >= %{version}-%{release}
	Requires: nvidia-l4t-multimedia-utils >= %{version}-%{release}
	Requires: libbsd

	%description gstreamer
		gstreamer drivers for Jetson nano

%package init
	Summary: Nvidia L4T init
	Requires: nvidia-l4t-core >= %{version}-%{release}

	%description init
		init for Jetson nano

%package initrd
	Summary: Nvidia L4T initrd
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: nvidia-l4t-xusb-firmware >= %{version}-%{release}

	%description initrd
		initrd for Jetson nano

%package jetson-io
	Summary: Nvidia L4T jetson-io drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: python3
	Requires: dtc
	# Requires: mount
	Requires: util-linux
	# Requires: nvidia-l4t-kernel >= %{version}-%{release}

	%description jetson-io
		jetson-io drivers for Jetson nano

%package multimedia
	Summary: Nvidia L4T multimedia drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: nvidia-l4t-3d-core >= %{version}-%{release}
	Requires: nvidia-l4t-cuda >= %{version}-%{release}
	Requires: nvidia-l4t-multimedia-utils >= %{version}-%{release}

	%description multimedia
		multimedia drivers for Jetson nano

%package multimedia-utils
	Summary: Nvidia L4T multimedia-utils drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: libdatrie1
	Requires: libfontconfig1
	Requires: harfbuzz
	Requires: pango
	Requires: pixman
	Requires: libXau
	Requires: libXrender
	Requires: zlib

	%description multimedia-utils
		multimedia-utils drivers for Jetson nano

%package oem-config
	Summary: Nvidia L4T oem-config
	Requires: nvidia-l4t-core >= %{version}-%{release}

	%description oem-config
		oem-config for Jetson nano

%package tools
	Summary: Nvidia L4T tools
	Requires: nvidia-l4t-core >= %{version}-%{release}

	%description tools
		tools for Jetson nano

%package wayland
	Summary: Nvidia L4T wayland drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}

	%description wayland
		wayland drivers for Jetson nano

%package weston
	Summary: Nvidia L4T weston drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: libaudit
	Requires: libdrm2
	Requires: libevdev
	Requires: libffi6
	Requires: mesa-libgbm
	Requires: libinput
	Requires: libjpeg-turbo
	Requires: pam-devel
	Requires: libpng
	Requires: systemd-libs
	Requires: libunwind

	%description weston
		weston drivers for Jetson nano

%package x11
	Summary: Nvidia L4T x11 drivers
	Requires: nvidia-l4t-core >= %{version}-%{release}
	Requires: nvidia-l4t-firmware >= %{version}-%{release}
	Requires: nvidia-l4t-init >= %{version}-%{release}

	%description x11
		x11 drivers for Jetson nano

%package xusb-firmware
	Summary: Nvidia L4T xusb-firmware
	Requires: nvidia-l4t-core >= %{version}-%{release}

	%description xusb-firmware
		xusb-firmware for Jetson nano

%prep
	rm -rf %{nv_dir}
	# Hold on. We don't want ALL of /etc.
	mkdir -p %buildroot/etc/systemd/system %buildroot/etc/X11
	mkdir -p %buildroot/usr/lib/firmware/ %buildroot/usr/lib64/systemd/
	wget %{url} -P %{nv_dir}
	cd %{nv_dir}
	tar xvf Tegra210_Linux_%{version}.%{release}_aarch64.tbz2
	rm Tegra210_Linux_%{version}.%{release}_aarch64.tbz2

%build
	cd %{nv_dir}
	tar xvf Linux_for_Tegra/nv_tegra/nvidia_drivers.tbz2
	tar xvf Linux_for_Tegra/nv_tegra/config.tbz2

%install
	sed -e 's_/usr/lib/aarch64-linux-gnu_/usr/lib64/aarch64-linux-gnu/_' -i %{nv_dir}/etc/nv_tegra_release
	sed -e 's_/usr/lib/_/usr/lib64/_' -i %{nv_dir}/etc/nv_tegra_release
	cp %{nv_dir}/etc/nv_tegra_release %buildroot/etc/nv_tegra_release
	cp -r %{nv_dir}/etc/ld.so.conf.d %buildroot/etc/ld.so.conf.d
	sed -e 's_/usr/lib/aarch64-linux-gnu_/usr/lib64/aarch64-linux-gnu/_' -i %{nv_dir}/etc/ld.so.conf.d/ld.so.conf
	echo "/usr/lib64/aarch64-linux-gnu/tegra-egl" > %buildroot/etc/ld.so.conf.d/nvidia-tegra.conf
	cp %{nv_dir}/etc/ld.so.conf.d/ld.so.conf %buildroot/etc/ld.so.conf.d/ld.so.conf

	cp %{nv_dir}/etc/systemd/nv* %buildroot/etc/systemd/
	cp -d %{nv_dir}/etc/systemd/system/nv*service %buildroot/etc/systemd/system/
	cp %{nv_dir}/etc/asound.conf.* %buildroot/etc/
	
	# Get the udev rules & xorg config.
	cp -r %{nv_dir}/etc/udev/* %buildroot/etc/udev/
	cp -r %{nv_dir}/etc/X11/xorg.conf %buildroot/etc/X11/

	# Copy usr/lib/aarch64-linux-gnu -> usr/lib64/aarch64-linux-gnu.
	cp -r %{nv_dir}/usr/lib/aarch64-linux-gnu/ %buildroot/usr/lib64/
	
	# Same for lib/firmware, lib/systemd.
	cp -r %{nv_dir}/lib/firmware/* %buildroot/usr/lib/firmware/
	cp -r %{nv_dir}/lib/systemd/* %buildroot/usr/lib64/systemd/

	# Pass through these 2 in usr/lib64.
	cp -r %{nv_dir}/usr/lib/xorg %buildroot/usr/lib64/
	cp -r %{nv_dir}/usr/lib/nvidia %buildroot/usr/lib/
	
	# These are OK as well...
	cp -r %{nv_dir}/usr/share %buildroot/usr/share/
	cp -r %{nv_dir}/usr/bin %buildroot/usr/bin/
	# move sbin -> bin
	cp -r %{nv_dir}/usr/sbin/ %buildroot/usr/
	# pass through
	cp -r %{nv_dir}/var/ %buildroot/var/
	cp -r %{nv_dir}/opt/ %buildroot/opt/ 

	[[ ! -e %buildroot/usr/lib/firmware/gm20b ]] && mkdir %buildroot/usr/lib/firmware/gm20b
	pushd %buildroot/usr/lib/firmware/gm20b > /dev/null 2>&1
                ln -sf "../tegra21x/acr_ucode.bin" "acr_ucode.bin"
                ln -sf "../tegra21x/gpmu_ucode.bin" "gpmu_ucode.bin"
                ln -sf "../tegra21x/gpmu_ucode_desc.bin" \
                                "gpmu_ucode_desc.bin"
                ln -sf "../tegra21x/gpmu_ucode_image.bin" \
                                "gpmu_ucode_image.bin"
                ln -sf "../tegra21x/gpu2cde.bin" \
                                "gpu2cde.bin"
                ln -sf "../tegra21x/NETB_img.bin" "NETB_img.bin"
                ln -sf "../tegra21x/fecs_sig.bin" "fecs_sig.bin"
                ln -sf "../tegra21x/pmu_sig.bin" "pmu_sig.bin"
                ln -sf "../tegra21x/pmu_bl.bin" "pmu_bl.bin"
                ln -sf "../tegra21x/fecs.bin" "fecs.bin"
                ln -sf "../tegra21x/gpccs.bin" "gpccs.bin"
                popd > /dev/null

				
	# Add a symlink for the Vulkan ICD.
	mkdir -p %buildroot/etc/vulkan/icd.d
	ln -s /usr/lib64/aarch64-linux-gnu/tegra/nvidia_icd.json %buildroot/etc/vulkan/icd.d/nvidia_icd.json
	
	# And another one for EGL.
	mkdir -p %buildroot/usr/share/glvnd/egl_vendor.d
	ln -s /usr/lib64/aarch64-linux-gnu/tegra-egl/nvidia.json %buildroot/usr/share/glvnd/egl_vendor.d/

	# Refresh old symlinks from /usr/lib/* to /usr/lib64/*
	ln -sfn /usr/lib64/aarch64-linux-gnu/tegra/libcuda.so %buildroot/usr/lib64/aarch64-linux-gnu/libcuda.so
	ln -sfn /usr/lib64/aarch64-linux-gnu/tegra/libcuda.so.1.1 %buildroot/usr/lib64/aarch64-linux-gnu/tegra/libcuda.so
	ln -sfn /usr/lib64/aarch64-linux-gnu/tegra/libnvbuf_utils.so.1.0.0 %buildroot/usr/lib64/aarch64-linux-gnu/tegra/libnvbuf_utils.so
	ln -sfn /usr/lib64/aarch64-linux-gnu/tegra/libnvbufsurface.so.1.0.0 %buildroot/usr/lib64/aarch64-linux-gnu/tegra/libnvbufsurface.so
	ln -sfn /usr/lib64/aarch64-linux-gnu/tegra/libnvbufsurftransform.so.1.0.0 %buildroot/usr/lib64/aarch64-linux-gnu/tegra/libnvbufsurftransform.so
	ln -sfn /usr/lib64/aarch64-linux-gnu/tegra/libnvid_mapper.so.1.0.0 %buildroot/usr/lib64/aarch64-linux-gnu/tegra/libnvid_mapper.so

%post 3d-core
%if [ "${ARCH}" = "arm64" ]; then
	echo "/usr/lib/aarch64-linux-gnu/tegra-egl" > \
		/usr/lib/aarch64-linux-gnu/tegra-egl/ld.so.conf
%endif
ldconfig

%post init
	%if ! -e "/opt/nvidia/l4t-packages/.nv-l4t-disable-boot-fw-update-in-preinstall"
		exit 0
	%endif
	
	cat << EOF > /etc/fstab
	# /etc/fstab: static file system information.
	#
	# These are the filesystems that are always mounted on boot, you can
	# override any of these by copying the appropriate line from this file into
	# /etc/fstab and tweaking it as you see fit.  See fstab(5).
	#
	# <file system> <mount point>             <type>          <options>                               <dump> <pass>
	/dev/root            /                     ext4           defaults                                     0 1
	EOF
	
	cat << EOF > /etc/hosts
	127.0.0.1 localhost
	EOF
	
	cat << EOF > /etc/hostname
	EOF
	
	cat <<EOF >> /etc/modules
	# bluedroid_pm, supporting module for bluetooth
	bluedroid_pm
	# modules for camera HAL
	nvhost_vi
	# nvgpu module
	nvgpu
	EOF
	
	cat <<EOF >> /etc/default/locale
	LANG=en_US.UTF-8
	EOF
	
	# Add crypto and trusty to system group
	groupadd -r crypto
	groupadd -r trusty
	
%post weston
%if grep "weston-launch" "/etc/group" == 1
	# create system group
	groupadd -r "weston-launch"
%endif


%files 3d-core
/etc/vulkan/icd.d/nvidia_icd.json
/usr/lib64/aarch64-linux-gnu/libvulkan.so.1.2
/usr/lib64/aarch64-linux-gnu/tegra/libGLX_nvidia.so.0
/usr/lib64/aarch64-linux-gnu/tegra/libnvidia-glsi.so.%{version}.%{release}
/usr/lib64/aarch64-linux-gnu/tegra/libnvvulkan-producer.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvidia-eglcore.so.%{version}.%{release}
/usr/lib64/aarch64-linux-gnu/tegra/libnvidia-glvkspirv.so.%{version}.%{release}        
/usr/lib64/aarch64-linux-gnu/tegra/libnvwinsys.so                     
/usr/lib64/aarch64-linux-gnu/tegra/libnvidia-fatbinaryloader.so.%{version}.%{release}  
/usr/lib64/aarch64-linux-gnu/tegra/libnvidia-ptxjitcompiler.so.%{version}.%{release}   
/usr/lib64/aarch64-linux-gnu/tegra/libvulkan.so.1.2.132               
/usr/lib64/aarch64-linux-gnu/tegra/libnvidia-glcore.so.%{version}.%{release}           
/usr/lib64/aarch64-linux-gnu/tegra/libnvidia-rmapi-tegra.so.%{version}.%{release}      
/usr/lib64/aarch64-linux-gnu/tegra/nvidia_icd.json
/usr/lib64/aarch64-linux-gnu/tegra-egl/ld.so.conf
/usr/lib64/aarch64-linux-gnu/tegra-egl/libGLESv1_CM_nvidia.so.1
/usr/lib64/aarch64-linux-gnu/tegra-egl/nvidia.json
/usr/lib64/aarch64-linux-gnu/tegra-egl/libEGL_nvidia.so.0
/usr/lib64/aarch64-linux-gnu/tegra-egl/libGLESv2_nvidia.so.2
/usr/lib64/xorg/modules/drivers/nvidia_drv.so
/usr/lib64/xorg/modules/extensions/libglxserver_nvidia.so
/usr/share/doc/nvidia-tegra/LICENSE.libvulkan1.gz
/usr/share/doc/nvidia-l4t-3d-core/*
/usr/share/glvnd/egl_vendor.d/10_nvidia.json

%files camera
/usr/lib64/aarch64-linux-gnu/tegra/libnvapputil.so                
/usr/lib64/aarch64-linux-gnu/tegra/libnvcameratools.so            
/usr/lib64/aarch64-linux-gnu/tegra/libnvcamv4l2.so              
/usr/lib64/aarch64-linux-gnu/tegra/libnvargus.so                  
/usr/lib64/aarch64-linux-gnu/tegra/libnvcamerautils.so            
/usr/lib64/aarch64-linux-gnu/tegra/libnveglstream_camconsumer.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvargus_socketclient.so     
/usr/lib64/aarch64-linux-gnu/tegra/libnvcam_imageencoder.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvodm_imager.so           
/usr/lib64/aarch64-linux-gnu/tegra/libnvargus_socketserver.so     
/usr/lib64/aarch64-linux-gnu/tegra/libnvscf.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvcamlog.so
/usr/sbin/nvargus-daemon
/usr/sbin/nvtunerd
/usr/share/doc/nvidia-l4t-camera/*
/usr/share/doc/nvidia-tegra/LICENSE.libnvargus
/usr/share/doc/nvidia-tegra/LICENSE.libnvcam_imageencoder
/etc/systemd/system/nvargus-daemon.service
/etc/systemd/system/multi-user.target.wants/nvargus-daemon.service
/var/nvidia/nvcam/*

%files core
/etc/nv_tegra_release
/usr/lib64/aarch64-linux-gnu/libdrm_nvdc.so
/usr/lib64/aarch64-linux-gnu/tegra/lib/libdrm.so.2                      
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_force.so                
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvos.so                     
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvcolorutil.so                
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_generic.so              
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvphsd.so                   
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvdc.so                       
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_gpucompute.so           
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvphs.so                    
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvddk_2d_v2.so                
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_graphics.so             
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvrm_gpu.so                 
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvddk_vic.so                  
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_il.so                   
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvrm_graphics.so            
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvfnet.so                     
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_spincircle.so           
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvrm.so                     
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvfnetstoredefog.so           
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_tbc.so                  
/usr/lib64/aarch64-linux-gnu/tegra/lib/libsensors.hal-client.nvs.so   
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvfnetstorehdfx.so            
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_ui.so                   
/usr/lib64/aarch64-linux-gnu/tegra/lib/libsensors_hal.nvs.so          
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgbm.so                      
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvidia-tls.so.%{version}.%{release}          
/usr/lib64/aarch64-linux-gnu/tegra/lib/libsensors.l4t.no_fusion.nvs.so
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_boot.so                 
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvimp.so                                                      
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvgov_camera.so               
/usr/lib64/aarch64-linux-gnu/tegra/lib/libnvll.so
/usr/sbin/nvphsd
/usr/share/doc/nvidia-tegra/L4T_End_User_License_Agreement.txt.gz
/usr/share/doc/nvidia-tegra/LICENSE.tegra_sensors    
/usr/share/doc/nvidia-tegra/LICENSE.minigbm
/usr/share/doc/nvidia-l4t-core/*

%files configs
/etc/systemd/nvweston.sh
/etc/systemd/system/nvweston.service
/etc/systemd/system/multi-user.target.wants/nvweston.service
/usr/share/backgrounds/*
/usr/share/doc/nvidia-l4t-configs/*
/usr/share/doc/procps/copyright.compliant
/usr/share/doc/udev/copyright.compliant.gz

%files cuda
/usr/lib64/aarch64-linux-gnu/libcuda.so
/usr/lib64/aarch64-linux-gnu/tegra/libcuda.so
/usr/lib64/aarch64-linux-gnu/tegra/libcuda.so.1.1
/usr/share/doc/nvidia-l4t-cuda/*

%files firmware
/etc/systemd/nvwifibt-pre.sh
/etc/systemd/nvwifibt.sh
/etc/systemd/system/nvwifibt.service
/lib/firmware/*
/lib/systemd/system/bluetooth.service.d/nv-bluetooth-service.conf
/usr/sbin/brcm_patchram_plus
/usr/share/doc/bluez/copyright.compliant.gz
/usr/share/doc/nvidia-l4t-firmware/*
/usr/share/doc/nvidia-tegra/LICENSE.brcm_patchram_plus.gz
/usr/share/doc/nvidia-tegra/LICENSE.cypress_wifibt.gz
/usr/share/doc/nvidia-tegra/LICENSE.realtek_8822ce_wifibt.gz

%files graphics-demos
/usr/share/doc/nvidia-l4t-graphics-demos/*
/usr/src/nvidia/graphics_demos/*

%files gstreamer
/usr/bin/gst-install
/usr/bin/nvgstcapture-1.0
/usr/bin/nvgstplayer-1.0
/usr/lib64/aarch64-linux-gnu/gstreamer-1.0/*
/usr/lib64/aarch64-linux-gnu/libgstnvegl-1.0.so.0
/usr/lib64/aarch64-linux-gnu/libgstnvexifmeta.so
/usr/lib64/aarch64-linux-gnu/libgstnvivameta.so
/usr/lib64/aarch64-linux-gnu/libnvsample_cudaprocess.so
/usr/lib64/aarch64-linux-gnu/tegra/libnveglstreamproducer.so
/usr/share/doc/nvidia-l4t-gstreamer/*
/usr/share/doc/nvidia-tegra/LICENSE.gst-nvvideo4linux2.gz
/usr/share/doc/nvidia-tegra/LICENSE.gst-openmax.gz
/usr/share/doc/nvidia-tegra/LICENSE.gstvideocuda
/usr/share/doc/nvidia-tegra/LICENSE.libgstnvvideosinks

%files init
/etc/asound.conf.tegrahda
/etc/asound.conf.tegrasndt210ref
/etc/default/rng-tools
/etc/enctune.conf
/etc/lightdm/lightdm.conf.d/50-nvidia.conf
/etc/modprobe.d/bcmdhd.conf
/etc/nv/nvfirstboot
/etc/nvidia-container-runtime/host-files-for-container.d/l4t.csv
/etc/nvphsd.conf
/etc/nvphsd_common.conf
/etc/sysctl.d/90-tegra-settings.conf
/etc/systemd/nv.sh
/etc/systemd/nvfb.sh
/etc/systemd/nvgetty.sh
/etc/systemd/nvmemwarning.sh
/etc/systemd/nvzramconfig.sh
/etc/systemd/sleep.conf
/etc/systemd/system/nv.service
/etc/systemd/system/nvfb.service
/etc/systemd/system/nvgetty.service
/etc/systemd/system/nvmemwarning.service
/etc/systemd/system/nvphs.service
/etc/systemd/system/nvs-service.service
/etc/systemd/system/nvzramconfig.service
/etc/udev/rules.d/90-alsa-asound-tegra.rules
/etc/udev/rules.d/91-xorg-conf-tegra.rules
/etc/udev/rules.d/92-hdmi-audio-tegra.rules
/etc/udev/rules.d/99-nv-l4t-usb-device-mode.rules
/etc/udev/rules.d/99-nv-l4t-usb-host-config.rules
/etc/udev/rules.d/99-nv-ufs-mount.rules
/etc/udev/rules.d/99-nv-wifibt.rules
/etc/udev/rules.d/99-tegra-devices.rules
/etc/udev/rules.d/99-tegra-mmc-ra.rules
/etc/wpa_supplicant.conf
/etc/xdg/autostart/nvchrome.desktop
/etc/xdg/autostart/nvchrome.sh
/opt/nvidia/l4t-bootloader-config/nv-l4t-bootloader-config.service
/opt/nvidia/l4t-bootloader-config/nv-l4t-bootloader-config.sh
/opt/nvidia/l4t-usb-device-mode/LICENSE.filesystem.img
/opt/nvidia/l4t-usb-device-mode/nv-l4t-usb-device-mode-config.sh
/opt/nvidia/l4t-usb-device-mode/nv-l4t-usb-device-mode-runtime-start.sh
/opt/nvidia/l4t-usb-device-mode/nv-l4t-usb-device-mode-runtime-stop.sh
/opt/nvidia/l4t-usb-device-mode/nv-l4t-usb-device-mode-runtime.service
/opt/nvidia/l4t-usb-device-mode/nv-l4t-usb-device-mode-start.sh
/opt/nvidia/l4t-usb-device-mode/nv-l4t-usb-device-mode-state-change.sh
/opt/nvidia/l4t-usb-device-mode/nv-l4t-usb-device-mode-stop.sh
/opt/nvidia/l4t-usb-device-mode/nv-l4t-usb-device-mode.service
/usr/bin/nvidia-bug-report-tegra.sh
/usr/sbin/nvphsd_setup.sh
/usr/sbin/nvs-service
/usr/sbin/nvsetprop
/usr/share/alsa/cards/tegra-hda.conf
/usr/share/alsa/cards/tegra-snd-t210r.conf
/usr/share/alsa/init/postinit/00-tegra.conf
/usr/share/alsa/init/postinit/01-tegra-rt565x.conf
/usr/share/doc/network-manager/copyright.compliant
/usr/share/doc/nvidia-l4t-init/*
/usr/share/doc/rng-tools/copyright
/usr/share/polkit-1/actions/com.nvidia.pkexec.nvpmodel.policy
/usr/share/polkit-1/actions/com.nvidia.pkexec.tegrastats.policy

# %file initrd

%files jetson-io
/opt/nvidia/jetson-io/*
/usr/share/doc/nvidia-l4t-jetson-io/*

%files multimedia
/usr/lib64/aarch64-linux-gnu/tegra/libnvmm_contentpipe.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvavp.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvbufsurface.so.1.0.0
/usr/lib64/aarch64-linux-gnu/tegra/libnvbufsurftransform.so.1.0.0
/usr/lib64/aarch64-linux-gnu/tegra/libnvdsbufferpool.so.1.0.0
/usr/lib64/aarch64-linux-gnu/tegra/libnveventlib.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvexif.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvid_mapper.so.1.0.0
/usr/lib64/aarch64-linux-gnu/tegra/libnvjpeg.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvmedia.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvmm.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvmmlite.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvmmlite_image.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvmmlite_utils.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvmmlite_video.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvmm_parser.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvmm_utils.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvofsdk.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvomx.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvomxilclient.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvosd.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvparser.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvtestresults.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvtnr.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvtracebuf.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvtvmr.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvv4l2.so
/usr/lib64/aarch64-linux-gnu/tegra/libnvv4lconvert.so
/usr/lib64/aarch64-linux-gnu/tegra/libtegrav4l2.so
/usr/lib64/aarch64-linux-gnu/tegra/libv4l2_nvvidconv.so
/usr/lib64/aarch64-linux-gnu/tegra/libv4l2_nvvideocodec.so
/usr/share/doc/nvidia-l4t-multimedia/*
/usr/share/doc/nvidia-tegra/LICENSE.libnveventlib
/usr/share/doc/nvidia-tegra/LICENSE.libnvtracebuf
/usr/share/doc/nvidia-tegra/LICENSE.libnvv4l2.gz
/usr/share/doc/nvidia-tegra/LICENSE.libnvv4lconvert.gz

%files multimedia-utils
/usr/lib64/aarch64-linux-gnu/tegra/libnvbuf_fdmap.so.1.0.0
/usr/lib64/aarch64-linux-gnu/tegra/libnvbuf_utils.so.1.0.0
/usr/share/doc/nvidia-l4t-multimedia-utils/*

%files oem-config
/etc/nv-oem-config.conf.t210
/etc/systemd/nv-oem-config-post.sh
/etc/systemd/nv-oem-config.sh
/lib/systemd/system/nv-oem-config-debconf@.service
/lib/systemd/system/nv-oem-config-gui.service
/lib/systemd/system/nv-oem-config.service
/lib/systemd/system/nv-oem-config.target
usr/lib/nvidia/license/nvlicense
usr/lib/nvidia/license/nvlicense-templates.sh
usr/lib/nvidia/resizefs/nvresizefs-query
usr/lib/nvidia/resizefs/nvresizefs.sh
usr/lib/nvidia/resizefs/nvresizefs.templates
usr/lib/ubiquity/plugins/nvlicense.py
usr/lib/ubiquity/plugins/nvresizefs.py
usr/sbin/nv-oem-config-firstboot
usr/share/doc/nvidia-l4t-oem-config/*

%files tools
/etc/nvpmodel/nvpmodel_t210.conf
/etc/nvpmodel/nvpmodel_t210_jetson-nano.conf
/etc/xdg/autostart/nvpmodel_indicator.desktop
/usr/bin/jetson_clocks
/usr/bin/tegrastats
/usr/sbin/l4t_payload_updater_t210
/usr/sbin/nvpmodel
/usr/share/doc/nvidia-l4t-tools/*
/usr/share/doc/nvidia-tegra/NVIDIA_Trademark_License_Addendum_SW.pdf.gz
/usr/share/nvpmodel_indicator/nvpmodel.py
/usr/share/nvpmodel_indicator/nvpmodel_indicator.py
/usr/share/nvpmodel_indicator/nv_logo.svg

%files wayland
/usr/lib64/aarch64-linux-gnu/tegra/libnvidia-egl-wayland.so
/usr/share/doc/nvidia-l4t-wayland/*
/usr/share/egl/egl_external_platform.d/nvidia_wayland.json

%files weston
/etc/xdg/weston/weston.ini
/usr/lib64/aarch64-linux-gnu/tegra/weston/*
/usr/share/doc/nvidia-l4t-weston/*
/usr/share/weston/*

%files x11
/etc/X11/xorg.conf
/etc/X11/xorg.conf.jetson_e
/usr/share/doc/nvidia-l4t-x11/changelog.Debian.gz
/usr/share/doc/nvidia-l4t-x11/copyright

%files xusb-firmware
/lib/firmware/tegra21x_xusb_firmware
/opt/ota_package/t21x/xusb_only_payload
/usr/share/doc/nvidia-l4t-xusb-firmware/*

%clean
rm -rf %{buildroot}