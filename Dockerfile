FROM golang:bookworm as builder

WORKDIR /build

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
    	ldconfig


# Verify ImageMagick installation and header files
RUN ls -la /usr/include/ImageMagick-7
RUN pkg-config --cflags --libs MagickWand

# Set environment variables for CGO
ENV PKG_CONFIG_PATH="/usr/lib/pkgconfig" \
    CGO_CFLAGS="`pkg-config --cflags MagickWand`" \
    CGO_LDFLAGS="`pkg-config --libs MagickWand`" \
    CGO_CFLAGS_ALLOW='-Xpreprocessor' \
    LD_LIBRARY_PATH="/usr/lib"

COPY . .

RUN go mod tidy && \
    go install -tags no_pkgconfig -v gopkg.in/gographics/imagick.v3/imagick && \
    go build -o app && \
    mv app  /opt/app

FROM debian:bookworm-slim

MAINTAINER buzzxu <downloadxu@163.com>



ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    fontconfig \
    libjpeg62-turbo \
    libpng16-16 \
    libtiff5 \
    libwebp \
    libgif7 \
    libx11-6 \
    libgomp1 && \
    ln -sv /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /data/images /app && \
    apt-get clean && \
    apt-get autoremove -y && \
    apt-get autoclean && \
    rm -rf /var/lib/apt/lists/*

# Copy application and necessary files
COPY --from=builder /usr/lib/libMagick* /usr/lib/
COPY --from=builder /usr/lib/ImageMagick /usr/lib/ImageMagick
COPY --from=builder /build/app /app/
COPY docker/conf.yml /app/
COPY docker/run.sh /app/
COPY docker/default.png /data/images/
COPY assets/msyh.ttf /app/

RUN ldconfig

# Set environment variables
ENV TZ=Asia/Shanghai \
    LANG=C.UTF-8 \
    LD_LIBRARY_PATH="/usr/lib"

WORKDIR /app

EXPOSE 3000

ENTRYPOINT ["/bin/bash","run.sh"]

