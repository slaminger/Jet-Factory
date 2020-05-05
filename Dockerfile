FROM ubuntu:18.04
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && apt-get install -y git p7zip dtrx wget lvm2 qemu dosfstools qemu-user-static arch-install-scripts multipath-tools
RUN mkdir -p /builder/