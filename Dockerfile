FROM ubuntu:latest
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y git tar wget p7zip unzip parted xz-utils dosfstools lvm2 qemu qemu-user-static arch-install-scripts
RUN mkdir -p /root/builder/