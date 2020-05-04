# Maintainer: Your Name <youremail@domain.com>
%global nv_version "r32_Release_v4.2"

Name:		nvidia-l4t
Version:	R32
Release:	4.2
License:	GPL
BuildArch:	noarch
Summary:	Nvidia L4T drivers
URL:		https://developer.nvidia.com/embedded/L4T/%{nv_version}/t210ref_release_aarch64/Tegra210_Linux_%{version}.%{release}_aarch64.tbz2
Requires:	cairo gstreamer1 pango libEGL libXext libX11 gstreamer1-plugins-base

%define nv_dir   %{nv_name}-%{version}

%description
	Nvidia L4T drivers
		
%package 
	Summary: Xorg server common files
	Group: User Interface/X
	Requires: pixman >= 0.30.0	
	Requires: xkeyboard-config xkbcomp

%description common
	Common files shared among all X servers.

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

%install
	# Hold on. We don't want ALL of /etc.
	mkdir -p %buildroot/etc/
	mkdir -p %buildroot/etc/systemd/system

	sed -e 's_/usr/lib/aarch64-linux-gnu_/usr/lib64/aarch64-linux-gnu/_' -i %{nv_dir}/etc/nv_tegra_release
	sed -e 's_/usr/lib/_/usr/lib64/_' -i %{nv_dir}/etc/nv_tegra_release
	cp %{nv_dir}/etc/nv_tegra_release %buildroot/etc/nv_tegra_release
	cp -r %{nv_dir}/etc/ld.so.conf.d %buildroot/etc/ld.so.conf.d
	echo "/usr/lib64/aarch64-linux-gnu/tegra" > %buildroot/etc/ld.so.conf.d/nvidia-tegra.conf
	echo "/usr/lib64/aarch64-linux-gnu/tegra-egl" > %buildroot/etc/ld.so.conf.d/ld.so.conf

	cp %{nv_dir}/etc/systemd/nv* %buildroot/etc/systemd/
	cp -d %{nv_dir}/etc/systemd/system/nv*service %buildroot/etc/systemd/system/
	cp %{nv_dir}/etc/asound.conf.* %buildroot/etc/
	
	# Get the udev rules & xorg config.
	cp -r %{nv_dir}/etc/udev/ %buildroot/etc/udev
	mkdir %buildroot/etc/X11
	cp -r %{nv_dir}/etc/X11/xorg.conf %buildroot/etc/X11/

	mkdir -p %buildroot/usr/lib/firmware/ %buildroot/usr/lib64/systemd/
	
	# Copy usr/lib/aarch64-linux-gnu -> usr/lib64/aarch64-linux-gnu.
	cp -r %{nv_dir}/usr/lib/aarch64-linux-gnu/ %buildroot/usr/lib64/
	
	# Same for lib/firmware, lib/systemd.
	cp -r %{nv_dir}/lib/firmware/* %buildroot/usr/lib/firmware/
	cp -r %{nv_dir}/lib/systemd/* %buildroot/usr/lib64/systemd/

	# Pass through these 2 in usr/lib64.
	cp -r %{nv_dir}/usr/lib/xorg %buildroot/usr/lib64/
	cp -r %{nv_dir}/usr/lib/nvidia %buildroot/usr/lib64/
	
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

%post
	/sbin/ldconfig

%files
/usr/*
/etc/*
/opt/*
/var/*

%clean
rm -rf %{buildroot}