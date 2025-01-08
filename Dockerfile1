# Build stage for ImageMagick
FROM golang:bookworm as builder

WORKDIR /build

# Install ImageMagick development packages first
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

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Install ImageMagick from source
ARG IMAGEMAGICK_VERSION=7.1.1-43
RUN cd /tmp && \
    wget https://www.imagemagick.org/download/ImageMagick-${IMAGEMAGICK_VERSION}.tar.gz && \
    tar -xvf ImageMagick-${IMAGEMAGICK_VERSION}.tar.gz && \
    cd ImageMagick-${IMAGEMAGICK_VERSION} && \
    ./configure --prefix=/usr \
        --enable-shared \
        --disable-static \
        --with-modules \
        --without-magick-plus-plus \
        --without-perl \
        --disable-openmp \
        --with-gvc=no \
        --disable-docs && \
    make -j$(nproc) && \
    make install && \
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

# Copy source code
COPY . .

# Build the application
RUN go mod tidy && \
    go install -tags no_pkgconfig gopkg.in/gographics/imagick.v3/imagick && \
    go build -o app

# Final stage
FROM debian:bookworm-slim

LABEL maintainer="buzzxu <downloadxu@163.com>"

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        fontconfig \
        libjpeg62-turbo \
        libpng16-16 \
        libtiff5 \
        libwebp7 \
        libgif7 \
        libx11-6 \
        libmagickwand-6.q16 && \
    rm -rf /var/lib/apt/lists/* && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /data/images /app

# Copy application and necessary files
COPY --from=builder /usr/lib/libMagick* /usr/lib/
COPY --from=builder /usr/lib/ImageMagick-7.1 /usr/lib/ImageMagick-7.1
COPY --from=builder /build/app /app/
COPY docker/conf.yml /app/
COPY docker/run.sh /app/
COPY docker/default.png /data/images/
COPY assets/msyh.ttf /app/

# Update library cache
RUN ldconfig

# Set environment variables
ENV TZ=Asia/Shanghai \
    LANG=C.UTF-8 \
    LD_LIBRARY_PATH="/usr/lib"

WORKDIR /app

EXPOSE 3000

ENTRYPOINT ["/bin/bash", "run.sh"]