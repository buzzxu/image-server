FROM golang:bookworm as build

WORKDIR $GOPATH/src/image-server
ADD . $GOPATH/src/image-server
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io

# Install build dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    wget \
    build-essential \
    pkg-config \
    libjpeg-dev \
    libpng-dev \
    libtiff-dev \
    libwebp-dev \
    libgif-dev \
    libx11-dev \
    libmagickwand-dev && \
    rm -rf /var/lib/apt/lists/*

RUN cd && \
        pkg-config --cflags --libs MagickWand && \
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
    	make -j$(nproc) && make install && \
    	ldconfig /usr/local/lib && \
    export CGO_CFLAGS="-I`pkg-config --cflags MagickWand`"; \
    export CGO_LDFLAGS="-I`pkg-config --libs MagickWand`"; \
    export CGO_CFLAGS_ALLOW='-Xpreprocessor'; \
    rm -rf $GOPATH/pkg/linux_amd64/gopkg.in/gographics/imagick.v3 && \
    cd $GOPATH/src/image-server && go install -tags no_pkgconfig -v gopkg.in/gographics/imagick.v3/imagick && \
    go build -o app && \
    mv app  /opt/app


FROM debian:bookworm-slim

MAINTAINER buzzxu <downloadxu@163.com>

WORKDIR /app
COPY --from=build /opt/app /app
COPY --from=build /root/ImageMagick.tar.gz /tmp/ImageMagick.tar.gz

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y wget build-essential pkg-config fontconfig libjemalloc-dev \
    libjpeg-dev libpng-dev libtiff-dev libwebp-dev \
    libgif-dev libx11-dev --no-install-recommends libmagickwand-dev && \
    cd  /tmp && \
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

