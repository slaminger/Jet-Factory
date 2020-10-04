FROM ubuntu:19.10

ARG DEBIAN_FRONTEND=noninteractive
RUN apt update -y && apt install -y \
	qemu \
	qemu-user-static \
	binfmt-support \
	arch-install-scripts \
	linux-image-generic \
	libguestfs-tools \
	wget \
	p7zip-full \
	xz-utils \
	zerofree \
	libarchive-tools \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /build
VOLUME /out

ARG DISTRO
ENV DISTRO=${DISTRO}
ARG DEVICE
ENV DEVICE=${DEVICE}
ARG HEKATE
ENV HEKATE=${HEKATE}
ARG HEKATE_ID
ENV HEKATE_ID=${HEKATE_ID}

COPY configs configs/
COPY src src/
COPY utils utils/

RUN find /build -type f -iname "*.sh" -exec chmod +x {} \;

CMD /build/src/entrypoint.sh /out/
