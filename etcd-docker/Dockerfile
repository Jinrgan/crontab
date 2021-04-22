FROM ubuntu:18.04

LABEL maintainer="Jinrgan <jinrgan@163.com>"

RUN sed -i "s@archive.ubuntu.com@mirrors.aliyun.com@g" /etc/apt/sources.list \
    && apt-get update
RUN apt-get -y install wget

WORKDIR /root
RUN wget https://github.com/etcd-io/etcd/releases/download/v3.3.9/etcd-v3.3.9-linux-amd64.tar.gz \
    && tar -xvzf etcd-v3.3.9-linux-amd64.tar.gz \
    && mv etcd-v3.3.9-linux-amd64/etcd /usr/local/bin/

CMD nohup etcd --listen-client-urls 'http://0.0.0.0:2379' --advertise-client-urls 'http://0.0.0.0:2379'
