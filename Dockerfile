FROM golang:alpine as build

WORKDIR $GOPATH/src/image-server
ADD . $GOPATH/src/image-server
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update && apk upgrade;  \
    apk --no-cache add pkgconf gcc g++ imagemagick-dev graphicsmagick-dev; \
    export CGO_CFLAGS="$(pkg-config --cflags MagickWand)" && export CGO_LDFLAGS="$(pkg-config --libs MagickWand)"; \
    cd $GOPATH/src/image-server && go install; \
    mv $GOPATH/bin/image-server /opt/app
#ENV IMAGEMAGICK_VERSION=7.0.8-68
#RUN apk update && apk upgrade;  \
#    apk --no-cache add axel zlib libpng-dev libjpeg-turbo-dev libwebp-dev giflib-dev libx11-dev xz-dev; \
#    axel -a -n 10 -o ImageMagick.tar.gz https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz; \
#    tar xvzf ImageMagick.tar.gz && \
#    cd ImageMagick* && \
#	./configure \
#	    --without-magick-plus-plus \
#	    --without-perl \
#	    --disable-openmp \
#	    --with-gvc=no \
#	    --disable-docs && \
#	make -j$(nproc) && make install && \
#	ldconfig /usr/local/lib && \
##	cd .. && rm -rf ImageMagick-${IMAGEMAGICK_VERSION}; \
#    pkg-config --cflags --libs MagickWand MagickCore;  \
#    go install -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick; \
##    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick; \
#    chmod a+x $GOPATH/bin/image-server;


FROM alpine

MAINTAINER buzzxu <downloadxu@163.com>

WORKDIR /app
COPY --from=build /opt/app /app
#
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk upgrade && \
    apk add --no-cache -U pkgconf zlib libpng libjpeg-turbo libwebp giflib tzdata xz-dev imagemagick && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    chmod a+x app && \
    mkdir -p /data/images && \
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*
#

ADD docker/conf.yml /app/conf.yml
ADD docker/run.sh /app/run.sh
ADD docker/default.png /data/images

ENV TZ Asia/Shanghai
ENV LANG C.UTF-8

EXPOSE 3000
ENTRYPOINT ["/bin/sh","run.sh"]

