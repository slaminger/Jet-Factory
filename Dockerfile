FROM ubuntu:19.10
ARG DEBIAN_FRONTEND=noninteractive

RUN apt update -y && apt upgrade -y && apt install -y software-properties-common
RUN add-apt-repository ppa:longsleep/golang-backports -y
RUN apt update -y && apt install -y qemu qemu-user-static linux-image-generic docker.io golang-go libguestfs-tools libguestfs-dev

WORKDIR /root/

ADD ./*.sh /root/
ADD ./*.go /root/
ADD ./go.* /root/
ADD ./*.json /root/

RUN chmod a+x ./*.sh
RUN go build

CMD /root/engine.sh
