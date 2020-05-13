FROM ubuntu:18.04

ARG DEBIAN_FRONTEND=noninteractive
RUN apt update -y && apt upgrade -y
RUN apt-get install -y git wget qemu qemu-user-static arch-install-scripts tar xz-utils unrar p7zip linux-image-generic

WORKDIR /tools
ADD ./*.sh /tools/
RUN chmod a+x /tools/*.sh