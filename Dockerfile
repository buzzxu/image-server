FROM golang:buster as build

WORKDIR $GOPATH/src/image-server
ADD . $GOPATH/src/image-server
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io

#RUN apt-get update && apt install -y apt-transport-https ca-certificates curl
#RUN echo \
#    deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free\
#    deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster-updates main contrib non-free\
#    deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster-backports main contrib non-free\
#    deb https://mirrors.tuna.tsinghua.edu.cn/debian-security buster/updates main contrib non-free\
#    > /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends

RUN apt install -y  -q libjpeg-dev libpng-dev libtiff-dev libwebp-dev libgif-dev libx11-dev libltdl-dev --no-install-recommends;

RUN cd && \
    	wget https://www.imagemagick.org/download/ImageMagick.tar.gz && \
    	tar -xvf ImageMagick.tar.gz && \
    	cd ImageMagick* && \
    	./configure --prefix=/usr \
            --enable-shared \
            --disable-static \
            --with-modules \
    	    --without-magick-plus-plus \
    	    --without-perl \
    	    --disable-openmp \
    	    --with-gvc=no \
    	    --disable-docs && \
    	pwd && \
    	make -j$(nproc) && make install && \
    	ldconfig

# Set environment variables for CGO
ENV PKG_CONFIG_PATH="/usr/lib/pkgconfig" \
    CGO_CFLAGS="`pkg-config --cflags MagickWand`" \
    CGO_LDFLAGS="`pkg-config --libs MagickWand`" \
    CGO_CFLAGS_ALLOW='-Xpreprocessor' \
    LD_LIBRARY_PATH="/usr/lib"

RUN rm -rf $GOPATH/pkg/linux_amd64/gopkg.in/gographics/imagick.v3; \
    cd $GOPATH/src/image-server && go install -tags no_pkgconfig -v gopkg.in/gographics/imagick.v3/imagick; \
    go build -o app; \
    mv app  /opt/app;


FROM debian:bullseye-slim

MAINTAINER buzzxu <downloadxu@163.com>

WORKDIR /app
COPY --from=build /opt/app /app
COPY --from=build /root/ImageMagick.tar.gz /tmp/ImageMagick.tar.gz

ENV DEBIAN_FRONTEND noninteractive

#RUN apt-get update && apt install -y apt-transport-https ca-certificates curl
#RUN echo \
#    deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free\
#    deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster-updates main contrib non-free\
#    deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster-backports main contrib non-free\
#    deb https://mirrors.tuna.tsinghua.edu.cn/debian-security buster/updates main contrib non-free\
#    > /etc/apt/sources.list

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y wget build-essential pkg-config fontconfig libjemalloc-dev \
    libjpeg-dev libpng-dev libtiff-dev libwebp-dev \
    libgif-dev libx11-dev --no-install-recommends && \
#    cd /tmp && \
#    wget https://github.com/jemalloc/jemalloc/releases/download/4.5.0/jemalloc-4.5.0.tar.bz2 && \
#    tar -xjvf jemalloc-4.5.0.tar.bz2 && \
#    cd jemalloc-4.5.0/ && \
#    ./configure --prefix=/usr/local/jemalloc && \
#    make -j$(nproc) && make install && \
#    echo /usr/local/jemalloc/lib >> /etc/ld.so.conf && \
#    ldconfig  && \
    cd  /tmp && \
#	wget https://www.imagemagick.org/download/ImageMagick.tar.gz && \
	tar -xvf ImageMagick.tar.gz && \
	cd ImageMagick* && \
	./configure --prefix=/usr \
	    --without-magick-plus-plus \
	    --without-perl \
	    --with-jemalloc \
	    --disable-openmp \
	    --with-gvc=no \
	    --disable-docs && \
	make -j$(nproc) && make install && \
	ldconfig /usr/local/lib && \
	rm /etc/localtime && \
    ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /data/images && \
    apt-get remove --purge -y wget build-essential pkg-config && \
    apt-get clean && \
    apt-get autoremove -y && \
    apt-get autoclean && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/*

ADD docker/conf.yml /app/conf.yml
ADD docker/run.sh /app/run.sh
ADD docker/default.png /data/images
ADD assets/msyh.ttf /app/msyh.ttf

ENV TZ Asia/Shanghai
ENV LANG C.UTF-8

EXPOSE 3000
ENTRYPOINT ["/bin/bash","run.sh"]


