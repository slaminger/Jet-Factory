FROM ubuntu:19.10
ARG DEBIAN_FRONTEND=noninteractive

RUN apt update -y 
RUN apt install -y qemu qemu-user-static \
				arch-install-scripts linux-image-generic \
				libguestfs-tools wget p7zip-full xz-utils

WORKDIR /root/
COPY configs configs/
COPY src src/
RUN chmod +x src/*.sh

VOLUME [ "/root/linux", "/root/android" ]
ENTRYPOINT [ "/root/src/entrypoint.sh" ]
