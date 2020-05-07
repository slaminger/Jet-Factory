FROM ubuntu:18.04
ARG DEBIAN_FRONTEND=noninteractive

RUN apt update -y && apt upgrade -y
RUN apt-get install -y git p7zip xz-utils unzip wget qemu dosfstools qemu-user-static arch-install-scripts

# Set default to arch
ENV DISTRO "arch"

WORKDIR /root
RUN chmod a+x ./*.go

CMD /root/jet-factory.go