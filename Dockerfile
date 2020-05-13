FROM ubuntu:18.04

ARG DEBIAN_FRONTEND=noninteractive
RUN apt update -y && apt upgrade -y
RUN apt-get install -y qemu qemu-user-static arch-install-scripts linux-image-generic

WORKDIR /root/
ADD ./ /root/
CMD ["./engine.sh"]