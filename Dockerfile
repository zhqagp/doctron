FROM golang:1.15.2-alpine as builder

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on

RUN mkdir -p /doctron
COPY . /doctron

RUN cd /doctron && \
    go build && \
    cp -fr doctron /usr/local/bin && \
    chmod +x /usr/local/bin/doctron

FROM lampnick/runtime:chromium-alpine

MAINTAINER lampnick <nick@lampnick.com>


# @build-example docker build . -f Dockerfile -t harbor.cpcti.com/library/doctron:latest



RUN set -ex \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && wget https://cpcti-log.oss-cn-beijing.aliyuncs.com/SourceHanSansCN.zip \
    && unzip SourceHanSansCN.zip \
    && cp SubsetOTF/CN/* /usr/share/fonts/ \
    && rm -rf SubsetOTF LICENSE.txt \
    && apk add --no-cache fontconfig ttf-dejavu \
    && fc-cache -fv

COPY --from=builder  /usr/local/bin/doctron /usr/local/bin/doctron
COPY conf/default.yaml /doctron.yaml
EXPOSE 8080
CMD ["dumb-init", "doctron", "--config", "/doctron.yaml"]


