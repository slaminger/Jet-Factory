FROM ubuntu:19.10
ARG DEBIAN_FRONTEND=noninteractive

RUN apt update -y && apt upgrade -y
RUN apt install -y software-properties-common
RUN add-apt-repository ppa:longsleep/golang-backports -y && apt update -y
RUN apt install -y qemu qemu-user-static arch-install-scripts linux-image-generic docker.io golang-go

ENV DISTRO="fedora"
WORKDIR /root/
ADD ./* /root/
RUN chmod a+x ./*.sh ./*.go
CMD /root/engine.sh