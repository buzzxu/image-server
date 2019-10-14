FROM golang:alpine as build

WORKDIR $GOPATH/src/image-server
COPY . $GOPATH/src/image-server

RUN go build


FROM alpine

RUN apk update && apk upgrade && \
    apk add --no-cache -U zlib libpng-dev libjpeg-turbo-dev libwebp-dev giflib-dev libx11-dev tzdata xz-dev&& \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /app && \
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

ENV IMAGEMAGICK_VERSION=7.0.8-68

RUN cd && \
    apk --no-cache add --virtual .build-dependencies wget curl ca-certificates gcc musl-dev  && \
	wget https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz && \
	tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
	cd ImageMagick* && \
	./configure \
	    --without-magick-plus-plus \
	    --without-perl \
	    --disable-openmp \
	    --with-gvc=no \
	    --enable-lzw \
	    --disable-docs && \
	make -j$(nproc) && make install && \
	ldconfig /usr/local/lib && \
	cd .. && rm -rf ImageMagick-${IMAGEMAGICK_VERSION}
	apk del .build-dependencies


WORKDIR /app
COPY . /app


ENV TZ Asia/Shanghai
ENV LANG C.UTF-8
ENV NODE_ENV production
ENTRYPOINT ["/bin/sh","run.sh"]

