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

RUN apt install -y  -q libjpeg-dev libpng-dev libtiff-dev libwebp-dev libgif-dev libx11-dev --no-install-recommends;

RUN cd && \
    	wget http://www.imagemagick.org/download/ImageMagick.tar.gz && \
    	tar -xvf ImageMagick.tar.gz && \
    	cd ImageMagick* && \
    	./configure --prefix=/usr \
    	    --without-magick-plus-plus \
    	    --without-perl \
    	    --disable-openmp \
    	    --with-gvc=no \
    	    --disable-docs && \
    	make -j$(nproc) && make install && \
    	ldconfig /usr/local/lib && \
    export CGO_CFLAGS="-I`pkg-config --cflags MagickWand`"; \
    export CGO_LDFLAGS="-I`pkg-config --libs MagickWand`"; \
    export CGO_CFLAGS_ALLOW='-Xpreprocessor'; \
#    export CGO_ENABLED=0 GOOS=linux GOARCH=amd64; \
#    export CGO_LDFLAGS="\
#    -Wl,-Bstatic \
#        `pkg-config --libs MagickWand MagickCore` \
#         -ljbig -ltiff -ljpeg -lwebp -llzma -lfftw3 -lbz2 -lgomp \
#    -Wl,-Bdynamic \
#        -llcms2 -llqr-1 -lglib-2.0 -lpng12 -lxml2 -lz -lm -ldl \
#    "; \
    rm -rf $GOPATH/pkg/linux_amd64/gopkg.in/gographics/imagick.v3; \
    cd $GOPATH/src/image-server && go install -tags no_pkgconfig -v gopkg.in/gographics/imagick.v3/imagick; \
    go build -o app; \
    mv app  /opt/app;


FROM debian:buster-slim

MAINTAINER buzzxu <downloadxu@163.com>

WORKDIR /app
COPY --from=build /opt/app /app

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
    apt-get install -y wget build-essential pkg-config \
    libjpeg-dev libpng-dev libtiff-dev libwebp-dev \
    libgif-dev libx11-dev && \
    cd && \
    wget https://github.com/jemalloc/jemalloc/releases/download/5.2.1/jemalloc-5.2.1.tar.bz2 && \
    tar -xjvf jemalloc-5.2.1.tar.bz2 && \
    cd jemalloc-5.2.1/ && \
    ./configure --prefix=/usr/local/jemalloc --libdir=/usr/local/lib && \
    make -j$(nproc) && make install && \
    rm -rf /tmp/jemalloc && \
    rm -rf jemalloc* && \
    cd  && \
	wget http://www.imagemagick.org/download/ImageMagick.tar.gz && \
	tar -xvf ImageMagick.tar.gz && \
	cd ImageMagick* && \
	./configure --prefix=/usr \
	    --without-magick-plus-plus \
	    --without-perl \
	    --with-jemalloc=yes \
	    --disable-openmp \
	    --with-gvc=no \
	    --disable-docs && \
	make -j$(nproc) && make install && \
	rm -rf ImageMagick* && \
	ldconfig /usr/local/lib && \
	rm /etc/localtime && \
    ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /data/images && \
    apt-get remove w--purge -y get build-essential pkg-config && \
    apt-get clean && \
    apt-get autoremove -y &&
    apt-get autoclean && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/*

ADD docker/conf.yml /app/conf.yml
ADD docker/run.sh /app/run.sh
ADD docker/default.png /data/images

ENV TZ Asia/Shanghai
ENV LANG C.UTF-8

EXPOSE 3000
ENTRYPOINT ["/bin/sh","run.sh"]

