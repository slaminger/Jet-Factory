FROM ubuntu:19.10
ARG DEBIAN_FRONTEND=noninteractive

RUN apt update -y 
RUN apt install -y qemu qemu-user-static \
				arch-install-scripts linux-image-generic \
				libguestfs-tools libguestfs-dev

VOLUME [ "/linux", "/android" ]
ENTRYPOINT [ "/builder/src/entrypoint.sh" ]
